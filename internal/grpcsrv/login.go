package grpcsrv

import (
	"context"
	"fmt"

	"github.com/game-sales-analytics/users-service/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s server) LoginWithEmail(ctx context.Context, req *pb.LoginWithEmailRequest) (*pb.LoginWithEmailReply, error) {
	fmt.Printf("%#v\n", req)
	s.logger.WithField("email", req.User.Email).Debug("handling login with email request")
	return nil, status.Error(codes.Canceled, "something has been cancelled")
}
