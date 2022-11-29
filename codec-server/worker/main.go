package main

import (
	"crypto/tls"
	"log"
	"os"

	codecserver "github.com/ktenzer/samples-go/codec-server"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	const clientCertPath string = "/home/ktenzer/temporal/certs/ca.pem"
	const clientKeyPath string = "/home/ktenzer/temporal/certs/ca.key"

	var c client.Client
	var err error
	var cert tls.Certificate

	if os.Getenv("MTLS") == "true" {

		cert, err = tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		if err != nil {
			log.Fatalln("Unable to load certs", err)
		}

		c, err = client.Dial(client.Options{
			HostPort:  os.Getenv("TEMPORAL_HOST_URL"),
			Namespace: os.Getenv("TEMPORAL_NAMESPACE"),
			ConnectionOptions: client.ConnectionOptions{
				TLS: &tls.Config{Certificates: []tls.Certificate{cert}},
			},
			DataConverter: codecserver.DataConverter,
		})
	} else {
		c, err = client.Dial(client.Options{
			HostPort:      os.Getenv("TEMPORAL_HOST_URL"),
			Namespace:     os.Getenv("TEMPORAL_NAMESPACE"),
			DataConverter: codecserver.DataConverter,
		})
	}

	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "codecserver", worker.Options{})

	w.RegisterWorkflow(codecserver.Workflow)
	w.RegisterActivity(codecserver.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
