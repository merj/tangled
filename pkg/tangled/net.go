package tangled

import (
	"fmt"
	"io"
	"net"
)

type NetTarget struct {
	network string
	address string
}

func NewNetTarget(network, address string) *NetTarget {
	return &NetTarget{
		network: network,
		address: address,
	}
}

func (t *NetTarget) Connect(clientReader io.Reader, clientWriter io.Writer) error {
	server, err := net.Dial(t.network, t.address)
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

type NetSource struct {
	network string
	address string
}

func NewNetSource(network, address string) *NetSource {
	return &NetSource{
		network: network,
		address: address,
	}
}

func (s *NetSource) ListenAndServe(t Target, l io.Writer) error {
	server, err := net.Listen(s.network, s.address)
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
