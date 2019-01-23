// Package api provides the local go-spacemesh API endpoints. e.g. json-http and grpc-http2
package main

import (
	"fmt"
	"strconv"

	//"github.com/spacemeshos/go-spacemesh/api/pb"
	"oracle_server/pb"
	//"github.com/spacemeshos/go-spacemesh/common"
	//"github.com/spacemeshos/go-spacemesh/log"
	//"strconv"

	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is a grpc server providing the Spacemesh api
type Server struct {
	wSwitch *WorldSwitch
	Server  *grpc.Server
	Port    uint
}

// Register
func (s Server) Register(ctx context.Context, in *pb.Registration) (*pb.SimpleMessage, error) {
	log.Printf("registering %v in world %v \n", in.ID, in.World)
	s.wSwitch.Get(in.World).Register(in.ID)
	return &pb.SimpleMessage{Value: "OK"}, nil
}

// Unregister
func (s Server) Unregister(ctx context.Context, in *pb.Registration) (*pb.SimpleMessage, error) {
	log.Printf("Unregistering %v in world %v \n", in.ID, in.World)
	s.wSwitch.Get(in.World).Unregister(in.ID)
	return &pb.SimpleMessage{Value: "OK"}, nil
}

// Echo returns the response for an echo api request
func (s Server) Validate(ctx context.Context, in *pb.ValidReq) (*pb.ValidRes, error) {
	log.Println("Got validate req")
	fmt.Println("inst : ", in.InstanceID)
	fmt.Println("committee ", in.CommitteeSize)
	fmt.Println("proof ", in.ID)

	v := &pb.ValidRes{Valid: s.wSwitch.Get(in.World).Validate(in.InstanceID, int(in.CommitteeSize), in.ID)}

	return v, nil
}

// Echo returns the response for an echo api request
func (s Server) ValidateMap(ctx context.Context, in *pb.ValidReq) (*pb.ValidList, error) {
	log.Println("Got validate req")
	fmt.Println("inst : ", in.InstanceID)
	fmt.Println("committee ", in.CommitteeSize)
	fmt.Println("proof ", in.ID)

	v := s.wSwitch.Get(in.World).ValidateMap(in.InstanceID, int(in.CommitteeSize), in.ID)

	return v, nil
}

// StopService stops the grpc service.
func (s Server) StopService() {
	s.Server.Stop()
}

// NewGrpcService create a new grpc service using config data.
func NewGrpcService(port int, ws *WorldSwitch) *Server {
	server := grpc.NewServer()
	return &Server{Server: server, Port: uint(port), wSwitch: ws}
}

// StartService starts the grpc service.
func (s Server) StartService(status chan bool) {
	go s.startServiceInternal(status)
}

// This is a blocking method designed to be called using a go routine
func (s Server) startServiceInternal(status chan bool) {
	addr := ":" + strconv.Itoa(int(s.Port))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("failed to listen", err)
		return
	}

	pb.RegisterServiceServer(s.Server, s)

	// SubscribeOnNewConnections reflection service on gRPC server
	reflection.Register(s.Server)

	log.Println("grpc API listening on port %d", s.Port)

	if status != nil {
		status <- true
	}

	// start serving - this blocks until err or server is stopped
	if err := s.Server.Serve(lis); err != nil {
		log.Println("grpc stopped serving", err)
	}

	if status != nil {
		status <- true
	}

}
