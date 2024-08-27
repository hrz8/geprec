package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	greeterv1 "github.com/hrz8/geprec/pkg/pb/greeter/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type server struct {
	greeterv1.UnimplementedGreeterServiceServer
}

func (s *server) SayHello(ctx context.Context, r *greeterv1.SayHelloRequest) (*greeterv1.SayHelloResponse, error) {
	return &greeterv1.SayHelloResponse{Message: "Hello " + r.Name}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":3009")
	if err != nil {
		log.Fatal("Cannot listen on tcp", err)
	}

	greeterv1.RegisterGreeterServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	cli, err := grpc.NewClient(":3009", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Cannot create grpc client", err)
	}

	mux := runtime.NewServeMux()
	ctx := context.Background()
	greeterv1.RegisterGreeterServiceHandler(ctx, mux, cli)
	httpServer := &http.Server{
		Addr:    ":3008",
		Handler: mux,
	}

	go func() {
		fmt.Println("Listening grpc server on :3009")
		grpcServer.Serve(lis)
	}()

	go func() {
		fmt.Println("Listening http server on :3008")
		httpServer.ListenAndServe()
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGINT)

	<-ch
}
