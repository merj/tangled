package tangled

import (
	"io"
	"os"
)

type FileTarget struct {
	name string
}

func NewFileTarget(name string) *FileTarget {
	return &FileTarget{
		name: name,
	}
}

func (t *FileTarget) Connect(clientReader io.Reader, clientWriter io.Writer) error {
	server, err := os.Open(t.name)
	if err != nil {
		return err
	}
	defer server.Close()

	err = Relay(clientReader, clientWriter, server, server)
	if err != nil {
		return err
	}

	return nil
}
