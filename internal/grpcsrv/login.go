package grpcsrv

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/game-sales-analytics/users-service/internal/auth"
	"github.com/game-sales-analytics/users-service/internal/pb"
	"github.com/game-sales-analytics/users-service/internal/validate"
)

func (s server) LoginWithEmail(ctx context.Context, in *pb.LoginWithEmailRequest) (*pb.LoginWithEmailReply, error) {
	s.logger.Debug("handling login request")

	form := validate.LoginForm{
		Email:           in.User.Email,
		Password:        in.User.Password,
		DeviceUserAgent: in.User.ContextualInfo.DeviceUserAgent,
		IP:              in.User.ContextualInfo.Ip,
	}
	if err := s.validator.ValidateLoginForm(ctx, form); nil != err {
		var validationErr *validate.ValidationError
		if errors.As(err, &validationErr) {
			return nil, status.Errorf(codes.InvalidArgument, `{"field":"%s","error":"%s"}`, validationErr.Field, validationErr.Message)
		}

		s.logger.WithError(err).WithField("err_code", "E_VALIDATE_LOGIN_FORM").Error("failed validating login form")
		return nil, errorInternal
	}

	creds := auth.LoginWithEmailCreds{
		Email:    in.User.Email,
		Password: in.User.Password,
		LoginDefaultCreds: auth.LoginDefaultCreds{
			UserIPAddress:       in.User.ContextualInfo.Ip,
			UserDeviceUserAgent: in.User.ContextualInfo.DeviceUserAgent,
		},
	}
	loginRes, err := s.auth.LoginWithEmail(ctx, creds)
	if nil != err {
		if errors.Is(err, auth.ErrUnauthenticated) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}

		s.logger.WithError(err).WithField("err_code", "E_LOGIN_WITH_EMAIL").Error("failed logging user in by email")
		return nil, errorInternal
	}

	return &pb.LoginWithEmailReply{
		AuthToken: &pb.LoginWithEmailReply_AuthToken{
			Id:                 loginRes.Token.ID,
			Token:              loginRes.Token.Value,
			NotBeforeDateTime:  timestamppb.New(loginRes.Token.NotBeforeDateTime),
			ExpirationDateTime: timestamppb.New(loginRes.Token.ExpirationDateTime),
		},
	}, nil
}
