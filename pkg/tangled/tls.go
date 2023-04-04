package tangled

import (
	"crypto/tls"
	"fmt"
	"io"
)

type TLSTarget struct {
	network string
	address string
	config  *tls.Config
}

func NewTLSTarget(network, address string, config *tls.Config) *TLSTarget {
	return &TLSTarget{
		network: network,
		address: address,
		config:  config,
	}
}

func (t *TLSTarget) Connect(clientReader io.Reader, clientWriter io.Writer) error {
	target, err := tls.Dial(t.network, t.address, t.config)
	if err != nil {
		return err
	}
	defer target.Close()

	err = Relay(clientReader, clientWriter, target, target)
	if err != nil {
		return err
	}

	return nil
}

type TLSSource struct {
	network string
	address string
	config  *tls.Config
}

func NewTLSSource(network, address string, config *tls.Config) *TLSSource {
	return &TLSSource{
		network: network,
		address: address,
		config:  config,
	}
}

func (s *TLSSource) ListenAndServe(t Target, l io.Writer) error {
	server, err := tls.Listen(s.network, s.address, s.config)
	if err != nil {
		return err
	}
	defer server.Close()

	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Fprintf(l, "%s\n", err)
			continue
		}

		go func() {
			defer client.Close()
			err := t.Connect(client, client)
			if err != nil {
				fmt.Fprintf(l, "%s\n", err)
			}
		}()
	}

	return nil
}
