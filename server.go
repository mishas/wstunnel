package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
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

func getTlsConfig() (*tls.Config, error) {
	tlscfg := new(tls.Config)
	tlscfg.ClientCAs = x509.NewCertPool()
	if ca, err := ioutil.ReadFile(path.Join(*certsDir, "cacert.pem")); err == nil {
		tlscfg.ClientCAs.AppendCertsFromPEM(ca)
	} else {
		return nil, fmt.Errorf("Failed reading CA certificate: %v", err)
	}

	if cert, err := tls.LoadX509KeyPair(path.Join(*certsDir, "/cert.pem"), path.Join(*certsDir, "/key.pem")); err == nil {
		tlscfg.Certificates = append(tlscfg.Certificates, cert)
	} else {
		return nil, fmt.Errorf("Failed reading client certificate: %v", err)
	}

	tlscfg.ClientAuth = tls.RequireAndVerifyClientCert

	return tlscfg, nil
}

func main() {
	flag.Parse()

	socks, err := socks5.New(&socks5.Config{})
	if err != nil {
		panic(err)
	}

	http.Handle("/", websocket.Handler(func(conn *websocket.Conn) { socks.ServeConn(conn) }))

	server := &http.Server{Addr: getAddr()}

	if *certsDir == "" {
		err = server.ListenAndServe()
	} else {
		if server.TLSConfig, err = getTlsConfig(); err != nil {
			panic(err)
		}
		err = server.ListenAndServeTLS("", "")
	}
	panic(err)
}
