package main

import (
	proto "ChitChat/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("172.20.10.2:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewChitChatClient(conn)

	Streamer, err := client.Chat(context.Background())
	name, _ := os.Hostname()
	Streamer.Send(&proto.ChatIn{Sender: name, Text: "Joined"})

	if err != nil {
		log.Fatalf("Not working")
	}
	go func() {
		for {
			msg, err := Streamer.Recv()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(msg.Sender, "Said: ")
			fmt.Println(">", msg.Text)
		}
	}()
	for {
		name, _ := os.Hostname()
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		Streamer.Send(&proto.ChatIn{Sender: name, Text: line})
	}
}
