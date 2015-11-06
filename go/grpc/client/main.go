package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/harukasan/testing/go/grpc/echo"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:11111", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewEchoClient(conn)

	for {
		stream, err := client.Stream(context.Background(), &pb.Message{"Hello, grpc"})
		if err != nil {
			log.Fatal(err)
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Println(err)
				break
			} else {
				log.Println(msg)
			}
		}
	}
}
