package main

import (
	"flag"
	"fmt"
	"net/http"
	"path"

	socks5 "github.com/armon/go-socks5"
	"golang.org/x/net/websocket"
)

var (
	certsDir = flag.String("certs_dir", "", "Directory of certs for starting a wss:// server, or empty for ws:// server. Expected files are: cert.pem and key.pem.")
	port     = flag.Int("port", 0, "The port to listen to, or empty for default (80 or 443).")
)

func getAddr() string {
	if *port == 0 {
		if *certsDir == "" {
			*port = 80
		} else {
			*port = 443
		}
	}

	return fmt.Sprintf(":%d", *port)
}

func main() {
	flag.Parse()

	server, err := socks5.New(&socks5.Config{})
	if err != nil {
		panic(err)
	}

	http.Handle("/", websocket.Handler(func(conn *websocket.Conn) { server.ServeConn(conn) }))

	if *certsDir == "" {
		if err := http.ListenAndServe(getAddr(), nil); err != nil {
			panic(err)
		}
	} else {
		cert := path.Join(*certsDir, "cert.pem")
		key := path.Join(*certsDir, "key.pem")

		if err := http.ListenAndServeTLS(getAddr(), cert, key, nil); err != nil {
			panic(err)
		}
	}
}
