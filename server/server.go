package main

import (
	proto "ChitChat/grpc"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type ChitChatDatabase struct {
	proto.UnimplementedChitChatServer
	clients map[string]*Client
}
type Client struct {
	name string
	id   string
	send chan *proto.ChatOut
}

func (s *ChitChatDatabase) addClient(c *Client) {
	fmt.Println(c.name, "has joined the Server")
	s.clients[c.id] = c
}

func (s *ChitChatDatabase) removeClient(c *Client) {
	fmt.Println(c.name, "has left the Server")
	delete(s.clients, c.id)
	close(c.send)
}

func NewServer() *ChitChatDatabase {
	return &ChitChatDatabase{clients: make(map[string]*Client)}
}

func (s *ChitChatDatabase) Chat(stream proto.ChitChat_ChatServer) error {
	clientId := fmt.Sprintf("%p", stream)
	chatIn, _ := stream.Recv()
	newClient := &Client{
		id:   clientId,
		name: chatIn.GetSender(),
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
		txt := strings.TrimSpace(in.Text)
		if txt == ".exit" {
			fmt.Println("Removing ", newClient.id)
			s.removeClient(newClient)
			return nil
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
	listener, err := net.Listen("tcp", "0.0.0.0:5050")
	if err != nil {
		log.Fatalf("Lorte program det virker ikke", err)
	}
	grpcServer := grpc.NewServer()
	svc := NewServer()
	fmt.Println("Server started at ", time.Now())

	proto.RegisterChitChatServer(grpcServer, svc)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("pis hamrende lorte pgram det virker ikke")
	}
}
