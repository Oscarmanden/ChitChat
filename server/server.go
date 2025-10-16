package main

import (
	proto "SimpleService/grpc"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ChitChatDatabase struct {
	proto.UnimplementedChitChatServer
}

func (c *ChitChatDatabase) JoinChat(ctx context.Context, req *proto.ParticipantName) (*proto.Join, error) {
	return &proto.Join{
		LogicalTime:     1324,
		ParticipantName: req.Join,
	}, nil
}

func (c *ChitChatDatabase) SendMessage(ctx context.Context, req *proto.RelevantChatInfo) (*proto.Chat, error) {
	return &proto.Chat{
		LogicalTime: 1324,
		Text:        req.Text,
		SenderName:  req.Username,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Lorte program det virker ikke")
	}
	grpcServer := grpc.NewServer()
	svc := &ChitChatDatabase{}
	proto.RegisterChitChatServer(grpcServer, svc)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("pis hamrende lorte pgram det virker ikke")
	}
}
