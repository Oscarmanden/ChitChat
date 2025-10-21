package main

import (
	proto "ChitChat/grpc"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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
	s.clients[c.id] = c
	msg := c.name + " " + "has joined the Server"
	logClientMessage("Client", c.id, "Join")
	out := &proto.ChatOut{
		Sender: "Server",
		Text:   msg,
		Ts:     time.Now().Unix(),
	}

	for _, c := range s.clients {
		select {
		case c.send <- out:
		}
	}
	logServerMessage("Server", "Broadcast")
}

func (s *ChitChatDatabase) removeClient(c *Client) {
	msg := c.name + " " + "has left the Server"
	logClientMessage("Client", c.id, "Left")
	out := &proto.ChatOut{
		Sender: "Server",
		Text:   msg,
		Ts:     time.Now().Unix(),
	}

	for _, c := range s.clients {
		select {
		case c.send <- out:
		}
	}
	logServerMessage("Server", "Broadcast")
	delete(s.clients, c.id)
	close(c.send)
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

	currClient := &Client{
		id:   clientId,
		name: chatIn.GetSender(),
		send: make(chan *proto.ChatOut, 32),
	}
	s.addClient(currClient)
	go func() {
		for msg := range currClient.send {
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
		txt := strings.TrimSpace(in.Text)
		if txt == ".exit" {
			s.removeClient(currClient)
			return nil
		}
		out := &proto.ChatOut{
			Sender: in.Sender,
			Text:   in.Text,
			Ts:     time.Now().Unix(),
			Ls:     serverLogicalTime,
		}
		logClientMessage("Client", currClient.id, "Message")
		for _, c := range s.clients {
			select {
			case c.send <- out:
			}
		}
	}
}

func logClientMessage(component string, clientId string, eventType string) {
	// [Client] ClientID; 025020502 [Joined] @ LS: 04:30:52
	fmt.Printf("[%s] ClientID; %s [%s] @ LS: %d\n",
		component, clientId, eventType, serverLogicalTime)
	ClockIncrement()
}

func logServerMessage(component string, eventType string) {
	// [Client] ClientID; 025020502 [Joined] @ LS: 04:30:52
	fmt.Printf("[%s] [%s] @ LS: %d\n",
		component, eventType, serverLogicalTime)
	ClockIncrement()
}

func main() {
	listener, err := net.Listen("tcp", "localhost:5050")
	if err != nil {
		log.Fatalf("Lorte program det virker ikke", err)
	}
	grpcServer := grpc.NewServer()
	svc := NewServer()
	logServerMessage("Server", "Started")

	proto.RegisterChitChatServer(grpcServer, svc)
	go func() {
		grpcServer.Serve(listener)
	}()
	for {
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		txt := strings.TrimSpace(line)
		if txt == ".shutdown" {
			grpcServer.Stop()
		}
		logServerMessage("Server", "Stopped")
		os.Exit(0)
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
