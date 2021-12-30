package grpcsrv

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/game-sales-analytics/users-service/internal/apm"
	"github.com/game-sales-analytics/users-service/internal/auth"
	"github.com/game-sales-analytics/users-service/internal/pb"
	"github.com/game-sales-analytics/users-service/internal/validate"
)

func (s server) Authenticate(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateReply, error) {
	s.logger.Debug("handling authenticate request")
	span := sentry.StartSpan(ctx, "authenticate", sentry.TransactionName("handle-authenticate-request"))
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

	form := validate.AuthenticateForm{
		Token: in.Token,
	}
	child := span.StartChild("validate-form")
	child.Status = sentry.SpanStatusOK
	if err := s.validator.ValidateAuthenticateForm(validate.NewContext(ctx, child), form); nil != err {
		defer child.Finish()

		var validationErr *validate.ValidationError
		if errors.As(err, &validationErr) {
			child.Status = sentry.SpanStatusInvalidArgument
			return nil, status.Errorf(codes.InvalidArgument, `{"field":"%s","error":"%s"}`, validationErr.Field, validationErr.Message)
		}

		child.Status = sentry.SpanStatusInternalError
		log := s.logger.WithError(err).WithField("err_code", "E_VALIDATE_AUTHENTICATE_FORM")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed validating authenticate form")
		return nil, errorInternal
	}
	child.Finish()

	child = span.StartChild("verify-token")
	child.Status = sentry.SpanStatusOK
	verificationResult, err := s.auth.VerifyToken(auth.NewContext(ctx, child), in.Token)
	if nil != err {
		defer child.Finish()

		if errors.Is(err, auth.ErrTokenNotVerified) || errors.Is(err, auth.ErrUserNotExists) {
			child.Status = sentry.SpanStatusUnauthenticated
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}

		child.Status = sentry.SpanStatusInternalError
		log := s.logger.WithError(err).WithField("err_code", "E_VERIFY_TOKEN")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed verifying user auth token")
		return nil, errorInternal
	}
	child.Finish()

	return &pb.AuthenticateReply{
		AuthenticatedUser: &pb.AuthenticateReply_AuthenticatedUser{
			Id:        verificationResult.User.ID,
			FirstName: verificationResult.User.FirstName,
			LastName:  verificationResult.User.LastName,
		},
	}, nil
}
