package main

import (
	"io"
	"log"
	"net"

	pb "github.com/AswinManojan/gRPC-demo/proto"
	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type helloServer struct {
	pb.GreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start the server %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &helloServer{})
	log.Printf("Server started at %v", lis.Addr())
	grpcServer.Serve(lis)
}

func (s *helloServer) SayHelloBidirectionalStreaming(stream pb.GreetService_SayHelloBidirectionalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		log.Printf("Got request with name : %v", req.Name)
		res := &pb.HelloResponse{
			Message: "Hello " + req.Name,
		}
		stream.Send(res)
	}
}
