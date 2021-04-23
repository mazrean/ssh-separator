package domain

import (
	"io"

	"github.com/mazrean/separated-webshell/domain/values"
)

type Connection struct {
	isTty      bool
	io         *values.ConnectionIO
	windowPipe chan *values.Window
}

func NewConnection(isTty bool, io *values.ConnectionIO) *Connection {
	return &Connection{
		isTty:      isTty,
		io:         io,
		windowPipe: make(chan *values.Window),
	}
}

func (c *Connection) IsTty() bool {
	return c.isTty
}

func (c *Connection) Stdin() io.Reader {
	return c.io.Stdin()
}

func (c *Connection) Stdout() io.Writer {
	return c.io.Stdout()
}

func (c *Connection) Stderr() io.Writer {
	return c.io.Stdout()
}

func (c *Connection) Close() error {
	return c.io.Close()
}

func (c *Connection) WindowSender() chan<- *values.Window {
	return c.windowPipe
}

func (c *Connection) WindowReceiver() <-chan *values.Window {
	return c.windowPipe
}
