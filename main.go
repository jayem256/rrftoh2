package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"rrftoh2/config"
	"rrftoh2/constants"
	"rrftoh2/h2server"
	"strconv"
	"time"

	"golang.org/x/net/http2"
)

func main() {
	var cfg config.Config

	// Read config json.
	err := config.Parse("./config.json", &cfg)

	if err != nil {
		panic(err)
	}

	addr := cfg.Server.BindAddr + constants.NETWORK_SEPARATOR + strconv.Itoa(cfg.Server.Port)
	res, err := net.ResolveTCPAddr(constants.NETWORK, addr)

	if err != nil {
		panic(err)
	}

	// Load cert.
	cer, err := tls.LoadX509KeyPair(cfg.Server.Cert, cfg.Server.Key)
	if err != nil {
		panic(err)
	}

	// Load client certs.
	acceptedCA := x509.NewCertPool()

	for _, fileName := range cfg.Client.CertPool {
		fb, err := os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		if !acceptedCA.AppendCertsFromPEM(fb) {
			fmt.Println("Could not add " + fileName + " to CA pool")
		}
	}

	// Listen for incoming connections.
	l, err := net.ListenTCP(constants.NETWORK, res)
	fmt.Println("Now listening for connections on " + addr)

	if err != nil {
		panic(err)
	}

	serve := &http2.Server{
		MaxUploadBufferPerConnection: int32(cfg.Server.WindowSize),
		MaxUploadBufferPerStream:     int32(cfg.Server.WindowSize),
		IdleTimeout:                  time.Second * 30,
	}

	// Wait for incoming connections.
	for {
		// Handle incoming connection.
		conn, err := l.AcceptTCP()
		if err == nil {
			// Set TCP_NODELAY.
			conn.SetNoDelay(true)
			// Handle TLS handshake and serve HTTP2.
			go h2server.ServeH2(conn, cer, acceptedCA, serve, &cfg)
		}
	}
}
