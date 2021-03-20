package ssh

import (
	"fmt"
	"log"

	"github.com/gliderlabs/ssh"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/service"
)

type SSH struct {
	*service.User
	*ssh.Server
}

func NewSSH(user *service.User) *SSH {
	server := ssh.Server{}
	server.PasswordHandler = func(ctx ssh.Context, password string) bool {
		isOK, err := user.SSHAuth(ctx, &domain.User{
			Name:     ctx.User(),
			Password: password,
		})
		if err != nil || !isOK {
			log.Printf("ssh login error: %+v\n", err)
			return false
		}

		return true
	}

	server.Handler = func(s ssh.Session) {
		_, winCh, isTty := s.Pty()
		newWinCh := make(chan *domain.Window)
		if isTty {
			go func(winCh <-chan ssh.Window, newWinCh chan<- *domain.Window) {
				for win := range winCh {
					newWinCh <- &domain.Window{
						Height: uint(win.Height),
						Width:  uint(win.Width),
					}
				}
			}(winCh, newWinCh)
		}
		err := user.SSHHandler(s.Context(), s.User(), isTty, newWinCh, s, s, s.Stderr())
		if err != nil {
			log.Fatalf("failed in ssh: %+v\n", err)
			return
		}
	}

	return &SSH{
		User:   user,
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
