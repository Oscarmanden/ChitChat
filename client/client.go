package main

import (
	proto "SimpleService/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	name, err := os.Hostname()

	client := proto.NewChitChatClient(conn)

	message, err := client.JoinChat(context.Background(), &proto.ParticipantName{Join: name})
	if err != nil {
		log.Fatalf("Not working " + name + "is faulty" + err.Error())
	}

	println(message.ParticipantName, message.LogicalTime)

	reader := bufio.NewReader(os.Stdin)
	for {
		chatMessage, _ := reader.ReadString('\n')
		req := &proto.RelevantChatInfo{
			Text:     chatMessage,
			Username: name,
		}
		resp, err := client.SendMessage(context.Background(), req)
		if err != nil {
			fmt.Println("SendMessage:", err)
			continue
		}
		fmt.Println(resp.SenderName, "said: >", resp.Text, "@", resp.LogicalTime)
	}
}
