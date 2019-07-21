package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"github.com/873314461/quic-file/client"
	"github.com/873314461/quic-file/server"
	"log"
	"math/big"
)

func main() {
	isServer := flag.Bool("s", false, "server mode")
	isClient := flag.Bool("c", false, "client mode")
	flag.Parse()

	if (*isServer && *isClient) || (!*isServer && !*isClient) {
		log.Fatalln("server or client?")
	}
	if *isServer {
		s := server.NewFileServer("[::]:8000", generateTLSConfig(), nil)
		s.Run()
	}
	if *isClient {
		client.Client("127.0.0.1:8000", "send.bin")
	}
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
