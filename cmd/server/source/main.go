package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"open-cluster-management.io/api/cloudevents/generic/options/mqtt"
	"open-cluster-management.io/api/test/integration/cloudevents/source"
)

const (
	clientCertFile = "/Users/liuwei/go/src/github.com/skeeey/mochimqtt/hack/client.pem"
	clientKeyFile  = "/Users/liuwei/go/src/github.com/skeeey/mochimqtt/hack/client-key.pem"

	serverCAFile = "/Users/liuwei/go/src/github.com/skeeey/mochimqtt/hack/root-ca.pem"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	mqttOptions := mqtt.NewMQTTOptions()
	mqttOptions.BrokerHost = "127.0.0.1:8883"
	mqttOptions.CAFile = serverCAFile
	mqttOptions.ClientCertFile = clientCertFile
	mqttOptions.ClientKeyFile = clientKeyFile

	_, err := source.StartResourceSourceClient(context.TODO(), mqttOptions)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
