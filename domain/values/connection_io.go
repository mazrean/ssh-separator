package values

import "io"

type ConnectionIO struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	closer func() error
}

func NewConnectionIO(stdin io.Reader, stdout io.Writer, stderr io.Writer, closer func() error) *ConnectionIO {
	return &ConnectionIO{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		closer: closer,
	}
}

func (t *ConnectionIO) Stdin() io.Reader {
	return t.stdin
}

func (t *ConnectionIO) Stdout() io.Writer {
	return t.stdout
}

func (t *ConnectionIO) Stderr() io.Writer {
	return t.stderr
}

func (t *ConnectionIO) Close() error {
	return t.closer()
}
