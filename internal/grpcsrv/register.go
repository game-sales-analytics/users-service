package grpcsrv

import (
	"context"
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/game-sales-analytics/users-service/internal/apm"
	"github.com/game-sales-analytics/users-service/internal/db/repository"
	"github.com/game-sales-analytics/users-service/internal/id"
	"github.com/game-sales-analytics/users-service/internal/passhash"
	"github.com/game-sales-analytics/users-service/internal/pb"
	"github.com/game-sales-analytics/users-service/internal/validate"
)

func (s server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	s.logger.Debug("handling registration request")
	span := sentry.StartSpan(ctx, "register", sentry.TransactionName("handle-register-request"))
	span.Status = sentry.SpanStatusOK
	defer span.Finish()

	traceID, err := apm.ReadOrGenerateTraceID(ctx)
	if nil != err {
		span.Status = sentry.SpanStatusFailedPrecondition

		log := s.logger.WithError(err).WithField("err_code", "E_READ_OT_GENERATE_TRACE_ID")
		apm.SetSpanTagsFromLogEntry(span, log)
		log.Error("failed read or generating trace id from context")

		return nil, errorInternal
	}
	span.TraceID = traceID

	form := validate.RegisterForm{
		Email:                in.Email,
		Password:             in.Password,
		PasswordConfirmation: in.PasswordConfirmation,
		FirstName:            in.FirstName,
		LastName:             in.LastName,
	}
	child := span.StartChild("validate-form")
	child.Status = sentry.SpanStatusOK
	normalizedForm, err := s.validator.ValidateRegisterForm(validate.NewContext(ctx, child), form)
	if nil != err {
		defer child.Finish()

		var validationErr *validate.ValidationError
		if errors.As(err, &validationErr) {
			child.Status = sentry.SpanStatusInvalidArgument
			return nil, status.Errorf(codes.InvalidArgument, `{"field":"%s","error":"%s"}`, validationErr.Field, validationErr.Message)
		}

		child.Status = sentry.SpanStatusInternalError
		log := s.logger.WithError(err).WithField("err_code", "E_VALIDATE_REGISTER_FORM")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed validating register form")
		return nil, errorInternal
	}
	child.Finish()

	child = span.StartChild("generate-user-id")
	child.Status = sentry.SpanStatusOK
	userID, err := id.GenerateUserID()
	if nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := s.logger.WithError(err).WithField("err_code", "E_GENERATE_USER_ID")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed generating ID for user")
		return nil, errorInternal
	}
	child.Finish()

	child = span.StartChild("hash-user-password")
	child.Status = sentry.SpanStatusOK
	hashedPasswd, err := passhash.HashPassword(in.Password)
	if nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := s.logger.WithError(err).WithField("err_code", "E_HASH_USER_PASSWORD")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed hashing user password")
		return nil, errorInternal
	}
	child.Finish()

	user := repository.NewUserToSave{
		ID:              userID,
		Email:           in.Email,
		NormalizedEmail: normalizedForm.Email,
		FirstName:       in.FirstName,
		LastName:        in.LastName,
		Password:        hashedPasswd,
		RegisteredAt:    time.Now(),
	}

	child = span.StartChild("save-user")
	child.Status = sentry.SpanStatusOK
	if err := s.repo.SaveNewUser(repository.NewDBOperationContext(ctx, child), user); nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := s.logger.WithError(err).WithField("err_code", "E_SAVE_USER")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed saving user")
		return nil, errorInternal
	}
	child.Finish()

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
