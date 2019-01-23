package main

import "log"

func main() {

	ws := NewWorldSwitch()

	// start api servers
	// start grpc if specified or if json rpc specified
	grpcAPIService := NewGrpcService(9595, ws)
	status := make(chan bool)
	grpcAPIService.StartService(status)
	<-status // first send says service started
	jsonAPIService := NewJSONHTTPServer(9595, 3030)
	jsonAPIService.StartService(status)
	<-status // first send says service started

	<-status // second send says service stopped
	<-status // first send says service stopped

	log.Println("Shutting down oracle server.")

}
