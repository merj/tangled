package tangled

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type WebsocketTarget struct {
	address  string
	protocol string
	origin   string
}

func NewWebsocketTarget(address, protocol, origin string) *WebsocketTarget {
	return &WebsocketTarget{
		address:  address,
		protocol: protocol,
		origin:   origin,
	}
}

func (t *WebsocketTarget) Connect(clientReader io.Reader, clientWriter io.Writer) error {
	server, err := websocket.Dial(t.address, t.protocol, t.origin)
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

type WebsocketSource struct {
	address string
}

func NewWebsocketSource(address string) *WebsocketSource {
	return &WebsocketSource{
		address: address,
	}
}

func (s *WebsocketSource) ListenAndServe(t Target, l io.Writer) error {
	relayHandler := func(client *websocket.Conn) {
		err := t.Connect(client, client)
		if err != nil {
			fmt.Fprintf(l, "%s\n", err)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(relayHandler))
	return http.ListenAndServe(s.address, mux)
}
