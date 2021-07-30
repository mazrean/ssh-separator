package ssh

import (
	"fmt"
	"log"

	"github.com/gliderlabs/ssh"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/mazrean/separated-webshell/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var connectionCounter = promauto.NewGauge(prometheus.GaugeOpts{
	Help:      "Number of connections",
	Namespace: "webshell",
	Name:      "ssh_connection_num",
})

type SSH struct {
	service.IUser
	service.IPipe
	*ssh.Server
}

func NewSSH(user service.IUser, pipe service.IPipe) *SSH {
	server := ssh.Server{}
	server.PasswordHandler = func(ctx ssh.Context, password string) bool {
		userName, err := values.NewUserName(ctx.User())
		if err != nil {
			return false
		}

		pw, err := values.NewPassword(password)
		if err != nil {
			return false
		}

		isOK, err := user.Auth(ctx, userName, pw)
		if err != nil || !isOK {
			log.Printf("ssh login error: %+v\n", err)
			return false
		}

		return true
	}

	server.Handler = func(s ssh.Session) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v\n", err)
			}
		}()
		_, winCh, isTty := s.Pty()
		tty := values.NewConnectionIO(s, s, s, s.Close)
		connection := domain.NewConnection(isTty, tty)
		newWinCh := connection.WindowSender()
		defer close(newWinCh)
		if isTty {
			go func(winCh <-chan ssh.Window, newWinCh chan<- *values.Window) {
				for win := range winCh {
					newWinCh <- values.NewWindow(uint(win.Height), uint(win.Width))
				}
			}(winCh, newWinCh)
		}

		userName, err := values.NewUserName(s.User())
		if err != nil {
			log.Printf("failed un UserName constructor: %+v", err)
			return
		}

		connectionCounter.Inc()
		defer connectionCounter.Dec()
		err = pipe.Pipe(s.Context(), userName, connection)
		if err != nil {
			log.Printf("failed in ssh: %+v\n", err)
			return
		}
	}

	return &SSH{
		IUser:  user,
		IPipe:  pipe,
		Server: &server,
	}
}

func (ssh *SSH) Start(port int) error {
	ssh.Server.Addr = fmt.Sprintf(":%d", port)
	err := ssh.Server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("listen and serve error: %w", err)
	}

	return nil
}
