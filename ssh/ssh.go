package ssh

import (
	"fmt"
	"log"

	"github.com/gliderlabs/ssh"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/mazrean/separated-webshell/service"
)

type SSH struct {
	service.IUser
	*ssh.Server
}

func NewSSH(user service.IUser) *SSH {
	server := ssh.Server{}
	server.PasswordHandler = func(ctx ssh.Context, password string) bool {
		isOK, err := user.SSHAuth(ctx, domain.NewUserWithPassword(domain.UserName(ctx.User()), domain.Password(password)))
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
		tty := values.NewTty(s, s, s)
		connection := domain.NewConnection(isTty, tty)
		newWinCh := connection.WindowSender()
		if isTty {
			go func(winCh <-chan ssh.Window, newWinCh chan<- *values.Window) {
				for win := range winCh {
					newWinCh <- values.NewWindow(uint(win.Height), uint(win.Width))
				}
			}(winCh, newWinCh)
		}
		err := user.SSHHandler(s.Context(), domain.UserName(s.User()), connection)
		if err != nil {
			log.Fatalf("failed in ssh: %+v\n", err)
			return
		}
	}

	return &SSH{
		IUser:  user,
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
