package grpcsrv

import (
	"context"

	"github.com/game-sales-analytics/users-service/internal/pb"
)

func (s server) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateReply, error) {
	return nil, nil
}
