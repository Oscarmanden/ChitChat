package main

import (
	proto "SimpleService/grpc"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewSimpleServiceClient(conn)

	message, err := client.GetSimpleMessage(context.Background(), &proto.HelloRequest{})
	if err != nil {
		log.Fatalf("Not working")
	}

	println(message.Reply)

}
