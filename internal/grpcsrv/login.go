package grpcsrv

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/game-sales-analytics/users-service/internal/apm"
	"github.com/game-sales-analytics/users-service/internal/auth"
	"github.com/game-sales-analytics/users-service/internal/pb"
	"github.com/game-sales-analytics/users-service/internal/validate"
)

func (s server) LoginWithEmail(ctx context.Context, in *pb.LoginWithEmailRequest) (*pb.LoginWithEmailReply, error) {
	s.logger.Debug("handling login request")
	span := sentry.StartSpan(ctx, "login-with-email", sentry.TransactionName("handle-login-with-email-request"))
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

	form := validate.LoginForm{
		Email:           in.Email,
		Password:        in.Password,
		DeviceUserAgent: in.DeviceUserAgent,
		IP:              in.Ip,
	}
	child := span.StartChild("validate-form")
	child.Status = sentry.SpanStatusOK
	if err := s.validator.ValidateLoginForm(validate.NewContext(ctx, child), form); nil != err {
		defer child.Finish()

		var validationErr *validate.ValidationError
		if errors.As(err, &validationErr) {
			child.Status = sentry.SpanStatusInvalidArgument
			return nil, status.Errorf(codes.InvalidArgument, `{"field":"%s","error":"%s"}`, validationErr.Field, validationErr.Message)
		}

		child.Status = sentry.SpanStatusInternalError
		log := s.logger.WithError(err).WithField("err_code", "E_VALIDATE_LOGIN_FORM")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed validating login form")
		return nil, errorInternal
	}
	child.Finish()

	creds := auth.LoginWithEmailCreds{
		Email:    in.Email,
		Password: in.Password,
		LoginDefaultCreds: auth.LoginDefaultCreds{
			UserIPAddress:       in.Ip,
			UserDeviceUserAgent: in.DeviceUserAgent,
		},
	}
	child = span.StartChild("auth-service-login-with-email")
	child.Status = sentry.SpanStatusOK
	loginCtx := auth.NewContext(ctx, child)
	loginRes, err := s.auth.LoginWithEmail(loginCtx, creds)
	if nil != err {
		defer child.Finish()

		if errors.Is(err, auth.ErrUnauthenticated) {
			child.Status = sentry.SpanStatusInvalidArgument
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}

		child.Status = sentry.SpanStatusInternalError
		log := s.logger.WithError(err).WithField("err_code", "E_LOGIN_WITH_EMAIL")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed logging user in by email")
		return nil, errorInternal
	}
	child.Finish()

	return &pb.LoginWithEmailReply{
		AuthToken: &pb.LoginWithEmailReply_AuthToken{
			Id:                 loginRes.Token.ID,
			Token:              loginRes.Token.Value,
			NotBeforeDateTime:  timestamppb.New(loginRes.Token.NotBeforeDateTime),
			ExpirationDateTime: timestamppb.New(loginRes.Token.ExpirationDateTime),
		},
	}, nil
}
