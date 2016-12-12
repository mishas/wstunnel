package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"path"

	"golang.org/x/net/proxy"
	"golang.org/x/net/websocket"
)

var (
	certsDir = flag.String("certs_dir", "", "Directory of certs for TLS connection to AMQP, or empty for non-TLS connection. "+
		"Expected files are: cacert.pem, cert.pem and key.pem.")
	serverName = flag.String("server_name", "", "Name of the server for TLS verification, or empty for default")

	targetHost = flag.String("target_host", "", "The target host:port to tunnel to")
	port       = flag.Int("port", 8080, "The local port to listen on")
)

func getTlsConfig() (*tls.Config, error) {
	if *certsDir == "" {
		return nil, nil
	}

	tlscfg := new(tls.Config)
	tlscfg.RootCAs = x509.NewCertPool()
	if ca, err := ioutil.ReadFile(path.Join(*certsDir, "cacert.pem")); err == nil {
		tlscfg.RootCAs.AppendCertsFromPEM(ca)
	} else {
		return nil, fmt.Errorf("Failed reading CA certificate: %v", err)
	}

	if cert, err := tls.LoadX509KeyPair(path.Join(*certsDir, "/cert.pem"), path.Join(*certsDir, "/key.pem")); err == nil {
		tlscfg.Certificates = append(tlscfg.Certificates, cert)
	} else {
		return nil, fmt.Errorf("Failed reading client certificate: %v", err)
	}

	if *serverName != "" {
		tlscfg.ServerName = *serverName
	}
	return tlscfg, nil
}

func getWsConfig() (*websocket.Config, error) {
	url := url.URL{Scheme: "ws", Host: *targetHost}
	if *certsDir != "" {
		url.Scheme = "wss"
	}

	config, err := websocket.NewConfig(url.String(), "http://localhost/")
	if err != nil {
		return nil, err
	}

	if config.TlsConfig, err = getTlsConfig(); err != nil {
		return nil, err
	}

	return config, nil
}

func iocopy(dst io.Writer, src io.Reader, c chan error) {
	_, err := io.Copy(dst, src)
	c <- err
}

type closeable interface {
	CloseWrite() error
}

func closeWrite(conn net.Conn) {
	if closeme, ok := conn.(closeable); ok {
		closeme.CloseWrite()
	}
}

func handleConnection(wsConfig *websocket.Config, conn net.Conn) {
	defer conn.Close()

	tcp, err := proxy.FromEnvironment().Dial("tcp", wsConfig.Location.Host)
	if err != nil {
		log.Print("proxy.FromEnvironment().Dial(): ", err)
		return
	}
	ws, err := websocket.NewClient(wsConfig, tls.Client(tcp, wsConfig.TlsConfig))
	if err != nil {
		log.Print("websocket.NewClient(): ", err)
		return
	}
	defer ws.Close()

	c := make(chan error, 2)
	go iocopy(ws, conn, c)
	go iocopy(conn, ws, c)

	for i := 0; i < 2; i++ {
		if err := <-c; err != nil {
			fmt.Print("io.Copy(): ", err)
			return
		}
		// If any of the sides closes the connection, we want to close the write channel.
		closeWrite(conn)
		closeWrite(tcp)
	}
}

func main() {
	flag.Parse()

	wsConfig, err := getWsConfig()
	if err != nil {
		panic(err)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("ln.Accept(): ", err)
			continue
		}
		go handleConnection(wsConfig, conn)
	}
}
