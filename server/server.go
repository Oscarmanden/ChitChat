package main

import (
	proto "ChitChat/grpc"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

var serverLogicalTime int64 = 0

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

func NewServer() *ChitChatDatabase {
	return &ChitChatDatabase{clients: make(map[string]*Client)}
}

func (s *ChitChatDatabase) Chat(stream proto.ChitChat_ChatServer) error {
	clientId := fmt.Sprintf("%p", stream)

	chatIn, _ := stream.Recv()

	// update server logical clock to highest value of own and received clock
	remoteClock := chatIn.GetLs()
	LogicalClockCompare(remoteClock)

	// increment server logical clock on recieve
	ClockIncrement()

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
		// increment before sending
		ClockIncrement()
		out := &proto.ChatOut{
			Sender: in.Sender,
			Text:   in.Text,
			Ts:     time.Now().Unix(),
			Ls:     serverLogicalTime,
		}

		for _, c := range s.clients {
			select {
			case c.send <- out:
			}
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:5050")
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

func ClockIncrement() {
	serverLogicalTime = serverLogicalTime + 1
}

func LogicalClockCompare(remoteClock int64) {
	serverLogicalTime = max(serverLogicalTime, remoteClock)
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
