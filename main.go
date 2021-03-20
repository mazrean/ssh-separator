package main

import (
	"os"
	"strconv"

	"github.com/mazrean/separated-webshell/repository"
)

func main() {
	strAPIPort := os.Getenv("API_PORT")
	strSSHPort := os.Getenv("SSH_PORT")

	apiPort, err := strconv.Atoi(strAPIPort)
	if err != nil {
		panic(err)
	}

	sshPort, err := strconv.Atoi(strSSHPort)
	if err != nil {
		panic(err)
	}

	db, err := repository.Setup()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server, err := InjectServer()
	if err != nil {
		panic(err)
	}

	api := server.API
	ssh := server.SSH

	go func() {
		panic(api.Start(apiPort))
	}()

	panic(ssh.Start(sshPort))
}
