package tangled

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type HTTPSSource struct {
	address  string
	certFile string
	keyFile  string
}

func NewHTTPSSource(address, certFile, keyFile string) *HTTPSSource {
	return &HTTPSSource{
		address:  address,
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (s *HTTPSSource) ListenAndServe(t Target, l io.Writer) error {
	relayHandler := func(client *websocket.Conn) {
		err := t.Connect(client, client)
		if err != nil {
			fmt.Fprintf(l, "%s\n", err)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/target", websocket.Handler(relayHandler))
	mux.Handle("/", http.FileServer(http.Dir("htdocs")))
	return http.ListenAndServeTLS(s.address, s.certFile, s.keyFile, mux)
}
