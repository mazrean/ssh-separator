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
}

func NewSSH(user *service.User) {
	ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		isOK, err := user.SSHAuth(ctx, &domain.User{
			Name:     ctx.User(),
			Password: password,
		})
		if err != nil || !isOK {
			return false
		}

		return true
	})

	ssh.Handle(func(s ssh.Session) {
		err := user.SSHHandler(s.Context(), s.User(), s, s, s)
		if err != nil {
			log.Fatalf("failed in ssh: %w", err)
			return
		}
	})
}

func (*SSH) Start(port int) error {
	err := ssh.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		return fmt.Errorf("listen and serve error: %w", err)
	}

	return nil
}
