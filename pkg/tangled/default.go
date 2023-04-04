package tangled

import (
	"crypto/tls"
	"fmt"
	"net/url"
)

func DefaultTarget(targetURL *url.URL) (Target, error) {
	switch targetURL.Scheme {
	case "file":
		if targetURL.Host != "" && targetURL.Host != "localhost" {
			return nil, fmt.Errorf("invalid host")
		}
		if targetURL.Path == "" {
			return nil, fmt.Errorf("invalid path")
		}
		return NewFileTarget(targetURL.Path), nil
	case "exec":
		if targetURL.Host != "" && targetURL.Host != "localhost" {
			return nil, fmt.Errorf("invalid host")
		}
		if targetURL.Path == "" {
			return nil, fmt.Errorf("invalid path")
		}
		return NewExecTarget(targetURL.Path), nil
	case "tcp":
		if targetURL.Host == "" {
			return nil, fmt.Errorf("invalid host")
		}
		return NewNetTarget("tcp", targetURL.Host), nil
	case "tls+tcp":
		if targetURL.Host == "" {
			return nil, fmt.Errorf("invalid host")
		}
		cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
		if err != nil {
			return nil, err
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		return NewTLSTarget("tcp", targetURL.Host, config), nil
	case "unix":
		if targetURL.Host != "" && targetURL.Host != "localhost" {
			return nil, fmt.Errorf("invalid host")
		}
		if targetURL.Path == "" {
			return nil, fmt.Errorf("invalid path")
		}
		return NewNetTarget("unix", targetURL.Path), nil
	case "ws":
		if targetURL.Host == "" {
			return nil, fmt.Errorf("invalid host")
		}
		return NewWebsocketTarget(targetURL.String(), "", targetURL.Host), nil
	default:
		return nil, fmt.Errorf("unknown scheme")
	}
}

func DefaultSource(sourceURL *url.URL) (Source, error) {
	switch sourceURL.Scheme {
	case "tcp":
		if sourceURL.Host == "" {
			return nil, fmt.Errorf("invalid host")
		}
		return NewNetSource("tcp", sourceURL.Host), nil
	case "tls+tcp":
		if sourceURL.Host == "" {
			return nil, fmt.Errorf("invalid host")
		}
		cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
		if err != nil {
			return nil, err
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		return NewTLSSource("tcp", sourceURL.Host, config), nil
	case "unix":
		if sourceURL.Host != "" || sourceURL.Path == "" {
			return nil, fmt.Errorf("invalid host")
		}
		return NewNetSource("unix", sourceURL.Path), nil
	case "ws":
		if sourceURL.Host == "" {
			return nil, fmt.Errorf("invalid host")
		}
		return NewWebsocketSource(sourceURL.Host), nil
	case "http":
		if sourceURL.Host == "" {
			return nil, fmt.Errorf("invalid host")
		}
		return NewHTTPSource(sourceURL.Host), nil
	case "https":
		if sourceURL.Host == "" {
			return nil, fmt.Errorf("invalid host")
		}
		return NewHTTPSSource(sourceURL.Host, "cert.pem", "key.pem"), nil
	default:
		return nil, fmt.Errorf("unknown scheme")
	}
}
