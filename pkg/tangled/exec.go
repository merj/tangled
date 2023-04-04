package tangled

import (
	"io"
	"os/exec"
)

type ExecTarget struct {
	command string
}

func NewExecTarget(command string) *ExecTarget {
	return &ExecTarget{
		command: command,
	}
}

func (t *ExecTarget) Connect(clientReader io.Reader, clientWriter io.Writer) error {
	server := exec.Command(t.command)

	serverReader, err := server.StdoutPipe()
	if err != nil {
		return err
	}

	server.Stderr = server.Stdout

	serverWriter, err := server.StdinPipe()
	if err != nil {
		return err
	}

	err = server.Start()
	if err != nil {
		return err
	}
	defer server.Process.Kill()

	err = Relay(clientReader, clientWriter, serverReader, serverWriter)
	if err != nil {
		return err
	}

	return server.Wait()
}
