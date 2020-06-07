package main

import (
	"flag"
	"fmt"
)

// listener needed data, normally should be in config or ENV, and injected from main
const (
	portDef     = "8585"
	endpointDef = "/save"
	protocolDef = "http"
	hostDef     = "localhost"
)

func main() {
	fileToSend := flag.String("file", "sample71kb.txt", "txt file to send")
	host := flag.String("host", hostDef, "host with listener")
	endpoint := flag.String("endpoint", endpointDef, "listener endpoint for file saving")
	protocol := flag.String("proto", protocolDef, "listener protocol")
	port := flag.String("port", portDef, "listener port")

	flag.Parse()

	url := createURL(*protocol, *host, *endpoint, *port)
	err := sender(*fileToSend, url)
	if err != nil {
		panic(err)
	}
}

// get listener url
func createURL(protocol, host, endpoint, port string) string {
	return fmt.Sprintf("%s://%s:%s%s", protocol, host, port, endpoint)
}
