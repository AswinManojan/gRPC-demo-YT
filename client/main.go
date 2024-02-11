package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/AswinManojan/gRPC-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8080"
)

func main() {
	conn, _ := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	client := pb.NewGreetServiceClient(conn)
	names := &pb.NamesList{
		Names: []string{"Aswin", "Manoj"},
	}
	callSayHelloBidirectionalStream(client, names)
}

func callSayHelloBidirectionalStream(client pb.GreetServiceClient, names *pb.NamesList) {
	stream, _ := client.SayHelloBidirectionalStreaming(context.Background())
	waitc := make(chan struct{})
	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			log.Println(message)
		}
		close(waitc)
	}()

	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		stream.Send(req)
		time.Sleep(2 * time.Second)
	}
	stream.CloseSend()
	<-waitc
}
