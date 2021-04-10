package values

import "io"

type ConnectionIO struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func NewConnectionIO(stdin io.Reader, stdout io.Writer, stderr io.Writer) *ConnectionIO {
	return &ConnectionIO{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
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
