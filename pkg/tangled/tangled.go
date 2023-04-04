package tangled

import (
	"io"
)

type Target interface {
	Connect(io.Reader, io.Writer) error
}

type Source interface {
	ListenAndServe(Target, io.Writer) error
}

func read(r io.Reader, ec chan error) chan []byte {
	rc := make(chan []byte)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				res := make([]byte, n)
				copy(res, buf[:n])
				rc <- res
			}
			if err != nil {
				if err == io.EOF {
					rc <- nil
				} else {
					ec <- err
				}
				break
			}
		}
	}()

	return rc
}

func Relay(clientReader io.Reader, clientWriter io.Writer, serverReader io.Reader, serverWriter io.Writer) error {
	errorChan := make(chan error)

	clientChan := read(clientReader, errorChan)
	serverChan := read(serverReader, errorChan)
	for {
		select {
		case data := <-clientChan:
			if data == nil {
				return nil
			}
			_, err := serverWriter.Write(data)
			if err != nil {
				return err
			}
		case data := <-serverChan:
			if data == nil {
				return nil
			}
			_, err := clientWriter.Write(data)
			if err != nil {
				return err
			}
		case err := <-errorChan:
			return err
		}
	}
}
