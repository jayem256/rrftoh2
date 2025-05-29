package h2server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"rrftoh2/config"
	"rrftoh2/constants"
	"rrftoh2/handler"

	"golang.org/x/net/http2"
)

// ServeH2 performs TLS handshake on connection and handles it as h2 connection
func ServeH2(conn net.Conn, cer tls.Certificate, clientCAs *x509.CertPool, server *http2.Server, cfg *config.Config) {
	// Wrap it as TLS connection.
	tlsconn := tls.Server(conn, &tls.Config{
		// Require client certificate.
		ClientAuth: tls.RequireAndVerifyClientCert,
		// Client cert pool.
		ClientCAs:  clientCAs,
		ServerName: cfg.Server.HostName,
		// Server cert.
		Certificates: []tls.Certificate{cer},
		/*
			Only include h2 (HTTP/2 over TLS).
			NOTE: Does not affect server hello message.
		*/
		NextProtos: []string{constants.HTTP2_ALPN},
		// Advertise TLS 1.3 as supported version.
		MinVersion: tls.VersionTLS13,
		// Check negotiated protocol.
		VerifyConnection: func(cs tls.ConnectionState) error {
			// Enforce h2 agreement.
			if cs.NegotiatedProtocol != constants.HTTP2_ALPN {
				// Returning error aborts the TLS connection.
				return errors.New(conn.RemoteAddr().String() +
					" did not agree to use " + constants.HTTP2_ALPN + " protocol")
			}
			return nil
		},
	})

	// Perform TLS handshake and pass to HTTP2 handler.
	if err := tlsconn.Handshake(); err == nil {
		server.ServeConn(tlsconn, &http2.ServeConnOpts{
			Handler: &handler.RequestHandler{
				Conn:        tlsconn,
				Hostname:    cfg.Server.HostName,
				DocRoot:     cfg.File.DocRoot,
				Compression: cfg.Compression.EnableGzip,
				Treshold:    cfg.Compression.CompressionTreshold,
				ReadSize:    cfg.File.BufferedRead,
				Omitted:     cfg.Compression.OmitExtensions,
				GzBuffer:    cfg.Compression.GzipBuffer,
			},
		})
	} else {
		fmt.Println(err.Error())
	}

}
