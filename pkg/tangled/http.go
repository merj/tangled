package tangled

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type HTTPSource struct {
	address string
}

func NewHTTPSource(address string) *HTTPSource {
	return &HTTPSource{
		address: address,
	}
}

func (s *HTTPSource) ListenAndServe(t Target, l io.Writer) error {
	relayHandler := func(client *websocket.Conn) {
		err := t.Connect(client, client)
		if err != nil {
			fmt.Fprintf(l, "%s\n", err)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/target", websocket.Handler(relayHandler))
	mux.Handle("/", http.FileServer(http.Dir("htdocs")))
	return http.ListenAndServe(s.address, mux)
}
