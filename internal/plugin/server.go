package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/krzysztofzaucha/protobuf-sandbox/internal"
	"github.com/krzysztofzaucha/protobuf-sandbox/internal/model"
	"github.com/krzysztofzaucha/protobuf-sandbox/internal/srv"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"math"
	"net"
	"time"
)

// Server is a server plugin symbol name.
var Server server

var errServerPlugin = errors.New("server")

type server struct {
	config *internal.Config
}

func (c *server) WithConfig(config *internal.Config) {
	c.config = config
}

// Execute method executes plugin logic.
func (c *server) Execute() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.config.Server.Host, c.config.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(math.MaxInt64),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Timeout: 5 * time.Second,
			},
		),
	}

	s := grpc.NewServer(opts...)
	model.RegisterPersonServiceServer(s, srv.New())

	fmt.Println("starting server...")
	fmt.Printf("hosting server on: %s\n", lis.Addr().String())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}
