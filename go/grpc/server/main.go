package main

import (
	"errors"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/harukasan/testing/go/grpc/echo"
)

type Server struct {
}

func (s *Server) Echo(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	return msg, nil
}

func (s *Server) Stream(msg *pb.Message, stream pb.Echo_StreamServer) error {
	for i := 1; ; i++ {
		if i%10 == 0 {
			return errors.New("hoge")
		}
		stream.Send(msg)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":11111")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	pb.RegisterEchoServer(server, &Server{})
	server.Serve(lis)
}
