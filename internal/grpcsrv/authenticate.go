package grpcsrv

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/game-sales-analytics/users-service/internal/auth"
	"github.com/game-sales-analytics/users-service/internal/pb"
	"github.com/game-sales-analytics/users-service/internal/validate"
)

func (s server) Authenticate(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateReply, error) {
	s.logger.Debug("handling authenticate request")

	form := validate.AuthenticateForm{
		Token: in.Token,
	}
	if err := s.validator.ValidateAuthenticateForm(ctx, form); nil != err {
		var validationErr *validate.ValidationError
		if errors.As(err, &validationErr) {
			return nil, status.Errorf(codes.InvalidArgument, `{"field":"%s","error":"%s"}`, validationErr.Field, validationErr.Message)
		}

		s.logger.WithError(err).WithField("err_code", "E_VALIDATE_AUTHENTICATE_FORM").Error("failed validating authenticate form")
		return nil, errorInternal
	}

	verificationResult, err := s.auth.VerifyToken(ctx, in.Token)
	if nil != err {
		if errors.Is(err, auth.ErrTokenNotVerified) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		if errors.Is(err, auth.ErrUserNotExists) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}

		s.logger.WithError(err).WithField("err_code", "E_VERIFY_TOKEN").Error("failed verifying user auth token")
		return nil, errorInternal
	}

	return &pb.AuthenticateReply{
		AuthenticatedUser: &pb.AuthenticateReply_AuthenticatedUser{
			Id:        verificationResult.User.ID,
			FirstName: verificationResult.User.FirstName,
			LastName:  verificationResult.User.LastName,
		},
	}, nil
}
