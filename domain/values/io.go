package values

import "io"

type IO struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func NewTty(stdin io.Reader, stdout io.Writer, stderr io.Writer) *IO {
	return &IO{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (t *IO) Stdin() io.Reader {
	return t.stdin
}

func (t *IO) Stdout() io.Writer {
	return t.stdout
}

func (t *IO) Stderr() io.Writer {
	return t.stderr
}
