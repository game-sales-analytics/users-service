package grpcsrv

import (
	"context"

	"github.com/game-sales-analytics/users-service/internal/pb"
)

func (s server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	return &pb.PingReply{
		Pong: true,
	}, nil
}
