package grpc

import (
	"net"

	"git.solsynth.dev/hydrogen/identity/pkg/grpc/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	proto.UnimplementedAuthServer
	proto.UnimplementedNotifyServer
}

func StartGrpc() error {
	listen, err := net.Listen("tcp", viper.GetString("grpc_bind"))
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	proto.RegisterAuthServer(server, &Server{})
	proto.RegisterNotifyServer(server, &Server{})

	reflection.Register(server)

	return server.Serve(listen)
}
