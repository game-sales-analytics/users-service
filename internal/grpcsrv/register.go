package grpcsrv

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/game-sales-analytics/users-service/internal/db/repository"
	"github.com/game-sales-analytics/users-service/internal/id"
	"github.com/game-sales-analytics/users-service/internal/passhash"
	"github.com/game-sales-analytics/users-service/internal/pb"
	"github.com/game-sales-analytics/users-service/internal/validate"
)

func (s server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	s.logger.Debug("handling registration request")

	form := validate.RegisterForm{
		Email:                in.Form.Email,
		Password:             in.Form.Password,
		PasswordConfirmation: in.Form.PasswordConfirmation,
		FirstName:            in.Form.FirstName,
		LastName:             in.Form.LastName,
	}
	normalizedForm, err := s.validator.ValidateRegisterForm(ctx, form)
	if nil != err {
		var validationErr *validate.ValidationError
		if errors.As(err, &validationErr) {
			return nil, status.Errorf(codes.InvalidArgument, `{"field":"%s","error":"%s"}`, validationErr.Field, validationErr.Message)
		}

		return nil, errorInternal
	}

	userID, err := id.GenerateUserID()
	if nil != err {
		s.logger.WithError(err).WithField("err_code", "E_GENERATE_USER_ID").Error("failed generating ID for user")
		return nil, errorInternal
	}

	hashedPasswd, err := passhash.HashPassword(in.Form.Password)
	if nil != err {
		s.logger.WithError(err).WithField("err_code", "E_HASH_USER_PASSWORD").Error("failed hashing user password")
		return nil, errorInternal
	}

	user := repository.NewUserToSave{
		ID:              userID,
		Email:           in.Form.Email,
		NormalizedEmail: normalizedForm.Email,
		FirstName:       in.Form.FirstName,
		LastName:        in.Form.LastName,
		Password:        hashedPasswd,
		RegisteredAt:    time.Now(),
	}
	if err := s.repo.SaveNewUser(ctx, user); nil != err {
		s.logger.WithError(err).WithField("err_code", "E_SAVE_USER").Error("failed saving user")
		return nil, errorInternal
	}

	return &pb.RegisterReply{
		RegisteredUser: &pb.RegisterReply_RegisteredUser{
			Id:           user.ID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        user.Email,
			RegisteredAt: timestamppb.New(user.RegisteredAt),
		},
	}, nil
}
