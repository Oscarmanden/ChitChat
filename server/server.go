package main

import (
	proto "ChitChat/grpc"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type ChitChatDatabase struct {
	proto.UnimplementedChitChatServer
	clients map[string]*Client
}
type Client struct {
	id   string
	send chan *proto.ChatOut
}

func (s *ChitChatDatabase) addClient(c *Client) {
	s.clients[c.id] = c
}

func NewServer() *ChitChatDatabase {
	return &ChitChatDatabase{clients: make(map[string]*Client)}
}

func (s *ChitChatDatabase) Chat(stream proto.ChitChat_ChatServer) error {
	clientId := fmt.Sprintf("%p", stream)
	newClient := &Client{
		id:   clientId,
		send: make(chan *proto.ChatOut, 32),
	}
	s.addClient(newClient)
	go func() {
		for msg := range newClient.send {
			if err := stream.Send(msg); err != nil {
				return
			}
		}
	}()

	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		out := &proto.ChatOut{
			Sender: in.Sender,
			Text:   in.Text,
			Ts:     time.Now().Unix(),
		}

		for _, c := range s.clients {
			select {
			case c.send <- out:
			}
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Lorte program det virker ikke")
	}
	grpcServer := grpc.NewServer()
	svc := NewServer()

	proto.RegisterChitChatServer(grpcServer, svc)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("pis hamrende lorte pgram det virker ikke")
	}
}
