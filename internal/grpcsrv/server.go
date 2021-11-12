package grpcsrv

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/game-sales-analytics/users-service/internal/pb"
	"github.com/sirupsen/logrus"
)

type GrpcService interface {
	Listen(host string, port uint) error
}

type server struct {
	pb.UnimplementedUsersServiceServer
	logger *logrus.Logger
}

func (s server) Listen(host string, port uint) error {
	s.logger.WithField("host", host).WithField("port", port).Debug("starting server")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if nil != err {
		s.logger.WithField("host", host).WithField("port", port).WithError(err).WithField("err_code", "E_SERVER_TCP_BIND").Error("failed to start listening at specified address")
		return err
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUsersServiceServer(grpcServer, server{})
	return grpcServer.Serve(lis)
}
