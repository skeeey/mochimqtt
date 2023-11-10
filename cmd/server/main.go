package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/skeeey/mochimqtt/pkg/server/broker"
)

func main() {
	serverCertPath := flag.String("server-cert-path", "/mochi-mqtt/server.pem", "The server cert path")
	serverKeyPath := flag.String("server-key-path", "/mochi-mqtt/server-key.pem", "The server key path")
	caPath := flag.String("ca-path", "/mochi-mqtt/root-ca.pem", "The ca path")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	serverCertPemBlock, err := os.ReadFile(*serverCertPath)
	if err != nil {
		log.Fatal(err)
	}

	serverKeyPemBlock, err := os.ReadFile(*serverKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.X509KeyPair(serverCertPemBlock, serverKeyPemBlock)
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	caPem, err := os.ReadFile(*caPath)
	if err != nil {
		log.Fatal(err)
	}

	if !certPool.AppendCertsFromPEM(caPem) {
		log.Fatal("failed to append client ca")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	server := mqtt.New(nil)

	// level := new(slog.LevelVar)
	// server.Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: level,
	// }))
	// level.Set(slog.LevelDebug)

	_ = server.AddHook(new(broker.CertAuthHook), nil)

	tcp := listeners.NewTCP("mqtt-tls", ":8883", &listeners.Config{
		TLSConfig: tlsConfig,
	})
	err = server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	server.Log.Warn("caught signal, stopping...")
	_ = server.Close()
	server.Log.Info("main.go finished")
}
