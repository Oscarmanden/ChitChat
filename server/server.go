package main

import (
	proto "SimpleService/grpc"
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type messageDatabaseServer struct {
	proto.UnimplementedSimpleServiceServer
	message string
}

func (m *messageDatabaseServer) GetSimpleMessage(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Reply: m.message,
	}, nil
}

func main() {
	server := &messageDatabaseServer{message: "ayo ayo"}
	tm := time.Now()
	server.message += " " + tm.Format("00:00:00")
	//server.message += "%s" +

	server.start_server()
}

func (s *messageDatabaseServer) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Lorte program det virker ikke")

	}

	proto.RegisterSimpleServiceServer(grpcServer, s)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("pis hamrende lorte program det virker ikke")
	}

}
