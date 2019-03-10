package main

import (
	"context"
	"fmt"
	"log"
	"net"

	server "github.com/alexmorten/spacy-server"
	"google.golang.org/grpc"
)

var gamePool *server.GamePool

func main() {
	gamePool = server.NewGamePool()
	s := &Server{
		Port: 4000,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	server.RegisterSpacyServerServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//Server for grpc
type Server struct {
	Port int
}

//GetCredentials for grpc-web
func (s *Server) GetCredentials(context.Context, *server.Empty) (*server.Credentials, error) {
	return gamePool.NewCredentials(), nil
}

//GetUpdates for grpc-web
func (s *Server) GetUpdates(credentials *server.Credentials, stream server.SpacyServer_GetUpdatesServer) error {
	streamWrap, err := gamePool.AddConnection(credentials, stream)
	if err != nil {
		return err
	}
	<-streamWrap.ShutdownC
	return nil
}

//Act for grpc-web
func (s *Server) Act(ctx context.Context, action *server.Action) (*server.Empty, error) {
	gamePool.HandleAction(action)
	return &server.Empty{}, nil
}
