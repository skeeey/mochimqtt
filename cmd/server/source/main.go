package main

import (
	"context"
	"fmt"
	"log"

	"k8s.io/klog/v2"

	"open-cluster-management.io/api/cloudevents/generic/options/mqtt"

	"github.com/skeeey/mochimqtt/pkg/server"
	"github.com/skeeey/mochimqtt/pkg/signal"
)

const (
	host = "mochi-mqtt-mochi-mqtt.apps.server-foundation-sno-r8b9r.dev04.red-chesterfield.com:443"

	serverCAFile = "/home/cloud-user/go/src/github.com/skeeey/mochimqtt/hack/certs/root-ca.pem"

	clientCertFile = "/home/cloud-user/go/src/github.com/skeeey/mochimqtt/hack/certs/source/client.pem"
	clientKeyFile  = "/home/cloud-user/go/src/github.com/skeeey/mochimqtt/hack/certs/source/client-key.pem"
)

func main() {
	shutdownCtx, cancel := context.WithCancel(context.TODO())
	shutdownHandler := signal.SetupSignalHandler()
	go func() {
		defer cancel()
		<-shutdownHandler
		klog.Infof("\nReceived SIGTERM or SIGINT signal, shutting down controller.\n")
	}()

	ctx, terminate := context.WithCancel(shutdownCtx)
	defer terminate()

	for i := 1; i < 1000; i++ {
		server.GetStore().Add(server.NewResource(fmt.Sprintf("cluster%d", i), "resource1"))
	}

	<-ctx.Done()
}

func startSource(ctx context.Context) {
	mqttOptions := mqtt.NewMQTTOptions()
	mqttOptions.BrokerHost = host
	mqttOptions.CAFile = serverCAFile
	mqttOptions.ClientCertFile = clientCertFile
	mqttOptions.ClientKeyFile = clientKeyFile

	_, err := server.StartResourceSourceClient(ctx, mqttOptions)
	if err != nil {
		log.Fatal(err)
	}
}
