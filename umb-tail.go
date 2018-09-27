package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

import stomp "github.com/go-stomp/stomp"

const (
	hostEndpoint = "messaging-devops-broker02.web.prod.ext.phx2.redhat.com:61612"
)

var (
	certFile = flag.String("cert", "someCertFile", "A PEM encoded certificate file.")
	keyFile  = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "/etc/pki/tls/certs/ca-bundle.crt", "A PEM encoded CA's certificate file.")
	queue    = flag.String("queue", "/queue/Consumer.client-rbean.go-test.VirtualTopic.eng.>", "A STOMP queue to subscribe to.")
)

//Connect to ActiveMQ and listen for messages
func main() {
	flag.Parse()

	// Load client cert
	log.Printf("Loading cert and key...")
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Load CA cert
	log.Printf("Loading CA cert...")
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatalln(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the TLS Connection next
	log.Printf("Creating TLS conn...")
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	netConn, err := tls.Dial("tcp", hostEndpoint, tlsConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer netConn.Close()

	// Now create the stomp connection
	log.Printf("Creating STOMP conn...")
	conn, err := stomp.Connect(netConn)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Disconnect()

	log.Printf("Subscribing...")
	sub, err := conn.Subscribe(*queue, stomp.AckAuto)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		log.Printf("Waiting for message...")
		msg := <-sub.C
		if msg == nil {
			log.Fatalln("received nil")
		}
		log.Printf(string(msg.Body))
	}

	err = sub.Unsubscribe()
	if err != nil {
		fmt.Println(err)
	}
}
