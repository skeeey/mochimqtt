package main

import (
	"context"

	"k8s.io/klog/v2"
	"open-cluster-management.io/api/cloudevents/generic/options/mqtt"
	"open-cluster-management.io/api/cloudevents/work"
	"open-cluster-management.io/api/cloudevents/work/agent/codec"

	"github.com/skeeey/mochimqtt/pkg/client"
	"github.com/skeeey/mochimqtt/pkg/signal"
)

const (
	serverCAFile = "/home/cloud-user/go/src/github.com/skeeey/mochimqtt/hack/certs/root-ca.pem"
	host         = "mochi-mqtt-mochi-mqtt.apps.server-foundation-sno-r8b9r.dev04.red-chesterfield.com:443"

	clusterName    = "cluster1"
	clientCertFile = "/home/cloud-user/go/src/github.com/skeeey/mochimqtt/hack/certs/cluster1/client.pem"
	clientKeyFile  = "/home/cloud-user/go/src/github.com/skeeey/mochimqtt/hack/certs/cluster1/client-key.pem"
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

	mqttOptions := mqtt.NewMQTTOptions()
	mqttOptions.BrokerHost = host
	mqttOptions.CAFile = serverCAFile
	mqttOptions.ClientCertFile = clientCertFile
	mqttOptions.ClientKeyFile = clientKeyFile

	clientHolder, err := work.NewClientHolderBuilder(clusterName, mqttOptions).
		WithClusterName(clusterName).
		WithCodecs(codec.NewManifestCodec(nil)).
		NewClientHolder(ctx)
	if err != nil {
		klog.Fatal(err)
	}

	informer := clientHolder.ManifestWorkInformer().Informer()

	controller := client.NewController(informer, clientHolder.ManifestWorks(clusterName))

	go informer.Run(ctx.Done())
	go controller.Run(1, ctx.Done())

	<-ctx.Done()
}
