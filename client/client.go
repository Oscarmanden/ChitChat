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
	Streamer.Send(&proto.ChatIn{Sender: name, Text: "Joined"})

	if err != nil {
		log.Fatalf("Not working")
	}
	messageBuffer := make([]*proto.ChatOut, 0)

	go func() {
		for {
			msg, err := Streamer.Recv()
			if err != nil {
				fmt.Println(err)
				return
			}
			// add msg to buffer
			fmt.Println("Received remoteTime to logicaltime ", msg.Ls, clientLogicalTime)
			messageBuffer = append(messageBuffer, msg)

			ClockIncrement()
			LogicalClockCompare(msg.Ls)

			sort.Slice(messageBuffer, func(i, j int) bool {
				fmt.Println("sorted buffer")
				return messageBuffer[i].Ls < messageBuffer[j].Ls

			})

			for _, msg := range messageBuffer {
				fmt.Println(msg.Sender, "Said: ")
				fmt.Println(">", msg.Text)
			}
			// clear buffer
			messageBuffer = messageBuffer[:0]
		}

	}()

	for {
		// increment before sending
		ClockIncrement()
		name, _ := os.Hostname()
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')

		Streamer.Send(&proto.ChatIn{Sender: name, Text: line, Ls: clientLogicalTime})
	}

}

func ClockIncrement() {

	clientLogicalTime = clientLogicalTime + 1
}

func LogicalClockCompare(remoteClock int64) {

	clientLogicalTime = max(clientLogicalTime, remoteClock)

}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
