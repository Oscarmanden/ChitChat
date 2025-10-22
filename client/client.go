package main

import (
	proto "ChitChat/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Logical Clock
var clientLogicalTime int64 = 0

func main() {

	conn, err := grpc.NewClient(":5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewChitChatClient(conn)
	Streamer, err := client.Chat(context.Background())

	name, _ := os.Hostname()
	onSend()
	Streamer.Send(&proto.ChatIn{Sender: name, Text: "Joined", Ls: clientLogicalTime})

	if err != nil {
		log.Fatalf("Not working")
	}
	messageBuffer := make([]*proto.ChatOut, 0)

	go func() {
		printedCount := 0
		for {
			msg, err := Streamer.Recv()
			if err != nil {
				fmt.Println(err)
				return
			}

			onRecieve(msg.Ls)
			// add msg to buffer
			messageBuffer = append(messageBuffer, msg)

			sort.Slice(messageBuffer, func(i, j int) bool {
				return messageBuffer[i].Ls < messageBuffer[j].Ls
			})

			for i := printedCount; i < len(messageBuffer); i++ {
				fmt.Println(msg.Sender, "Said: ")
				fmt.Println(">", msg.Text)
			}
			printedCount = len(messageBuffer)
		}
	}()

	for {
		name, _ := os.Hostname()
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		onSend()
		Streamer.Send(&proto.ChatIn{Sender: name, Text: line, Ls: clientLogicalTime})
	}

}

func onSend() {
	clientLogicalTime = clientLogicalTime + 1
}
func onRecieve(remote int64) {
	clientLogicalTime = max(clientLogicalTime, remote) + 1
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
