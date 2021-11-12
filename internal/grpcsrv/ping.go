package grpcsrv

import (
	"context"
	"math/rand"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/game-sales-analytics/users-service/internal/pb"
)

func (s server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	if rand.Float32() < 0.5 {
		return nil, status.Error(codes.Unavailable, "unable to process request for now. try later.")
	}
	return &pb.PingReply{
		Pong: true,
	}, nil
}
