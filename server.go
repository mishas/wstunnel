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
	certsDir  = flag.String("certs_dir", "", "Directory of certs for starting a wss:// server, or empty for ws:// server. Expected files are: cert.pem and key.pem.")
	httpPort  = flag.Int("http_port", 80, "The port to listen to for http responses")
	httpsPort = flag.Int("https_port", 443, "The port to listen to for https responses")
)

func getTlsConfig() (*tls.Config, error) {
	tlscfg := new(tls.Config)
	tlscfg.MinVersion = tls.VersionTLS12
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

func startServers(httpServer, httpsServer *http.Server) error {
	c := make(chan error)
	go func() { c <- httpServer.ListenAndServe() }()
	if httpsServer != nil {
		go func() { c <- httpsServer.ListenAndServeTLS("", "") }()
	}

	return <-c
}

func setDebugHandlers(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/generate_204", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("success\n")) })
	return mux
}

func main() {
	flag.Parse()

	socks, err := socks5.New(&socks5.Config{})
	if err != nil {
		panic(err)
	}

	httpMux := setDebugHandlers(http.NewServeMux())
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", *httpPort), Handler: httpMux}
	mainMux := httpMux

	var httpsServer *http.Server
	if *certsDir != "" {
		httpsMux := setDebugHandlers(http.NewServeMux())
		mainMux = httpsMux
		httpsServer = &http.Server{
			Addr: fmt.Sprintf(":%d", *httpsPort), Handler: httpsMux,
			// The next line disables HTTP/2, as this does not support websockets.
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}
		if httpsServer.TLSConfig, err = getTlsConfig(); err != nil {
			panic(err)
		}
	}

	mainMux.Handle("/", websocket.Handler(func(conn *websocket.Conn) { socks.ServeConn(conn) }))

	panic(startServers(httpServer, httpsServer))
}
