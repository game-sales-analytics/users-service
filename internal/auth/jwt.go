package auth

import (
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/game-sales-analytics/users-service/internal/apm"
)

type GeneratedToken struct {
	ID                 string
	Value              string
	NotBeforeDateTime  time.Time
	ExpirationDateTime time.Time
}

func (a *authsrv) generateToken(ctx Context, userID, secret string) (*GeneratedToken, error) {
	child := ctx.span.StartChild("generate-auth-token-id")
	tokenID, err := uuid.NewV4()
	if nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_GENERATE_JWT_TOKEN_ID")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed generating jwt token id")
		return nil, errors.New("unable to generate token id")
	}
	child.Finish()

	child = ctx.span.StartChild("set-iss-key")
	token := jwt.New()
	if err := token.Set(jwt.IssuerKey, "https://github.com/game-sales-analytics/users-service"); nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_ISSUER_CLAIM")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed setting jwt token issuer claim")
		return nil, errors.New("unable to set auth token issuer key")
	}
	child.Finish()

	child = ctx.span.StartChild("set-aud-key")
	if err := token.Set(jwt.AudienceKey, []string{"users"}); nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_AUDIENCE_CLAIM")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed setting jwt token audience claim")
		return nil, errors.New("unable to set auth token audience key")
	}
	child.Finish()

	child = ctx.span.StartChild("set-exp-key")
	exp := time.Now()
	if err := token.Set(jwt.ExpirationKey, exp.Add(time.Hour*24*7)); nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_EXPIRATION_CLAIM")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed setting jwt token expiration claim")
		return nil, errors.New("unable to set auth token expiration key")
	}
	child.Finish()

	child = ctx.span.StartChild("set-iat-key")
	if err := token.Set(jwt.IssuedAtKey, time.Now()); nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_ISSUED_AT_CLAIM")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed setting jwt token issued_at claim")
		return nil, errors.New("unable to set auth token issued_at key")
	}
	child.Finish()

	child = ctx.span.StartChild("set-jti-key")
	if err := token.Set(jwt.JwtIDKey, tokenID.String()); nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_JWT_ID_CLAIM")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed setting jwt token jwt_id claim")
		return nil, errors.New("unable to set auth token jwt_id key")
	}
	child.Finish()

	child = ctx.span.StartChild("set-nbf-key")
	nbf := time.Now()
	if err := token.Set(jwt.NotBeforeKey, nbf); nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_NOT_BEFORE_CLAIM")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed setting jwt token not_before claim")
		return nil, errors.New("unable to set auth token not_before key")
	}
	child.Finish()

	child = ctx.span.StartChild("set-sub-key")
	if err := token.Set(jwt.SubjectKey, userID); nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_SUBJECT_CLAIM")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed setting jwt token subject claim")
		return nil, errors.New("unable to set auth token subject key")
	}
	child.Finish()

	child = ctx.span.StartChild("sign-auth-token")
	opts := []jwt.SignOption{}
	serialized, err := jwt.Sign(token, jwa.HS512, []byte(secret), opts...)
	if nil != err {
		defer child.Finish()

		child.Status = sentry.SpanStatusInternalError
		log := a.logger.WithError(err).WithField("err_code", "E_SIGN_JWT_TOKEN")
		apm.SetSpanTagsFromLogEntry(child, log)
		log.Error("failed signing jwt token")
		return nil, errors.New("unable to sign auth token")
	}
	child.Finish()

	return &GeneratedToken{
		ID:                 tokenID.String(),
		Value:              string(serialized),
		NotBeforeDateTime:  nbf,
		ExpirationDateTime: exp,
	}, err
}

type tokenDecodeResult struct {
	userID string
}

func verifyToken(ctx Context, token, secret string) (*tokenDecodeResult, error) {
	parseOptions := []jwt.ParseOption{
		jwt.WithVerify(jwa.HS512, []byte(secret)),
		jwt.WithValidate(true),
		jwt.WithAudience("users"),
		jwt.WithIssuer("https://github.com/game-sales-analytics/users-service"),
		jwt.WithMinDelta(time.Second*10, jwt.ExpirationKey, jwt.IssuedAtKey),
	}
	span := ctx.span.StartChild("parse-raw-token")
	parsedToken, err := jwt.Parse([]byte(token), parseOptions...)
	if nil != err {
		defer span.Finish()

		span.Status = sentry.SpanStatusInvalidArgument
		return nil, errors.New("parsing and verifying token failed")
	}
	span.Finish()

	out := tokenDecodeResult{
		userID: parsedToken.Subject(),
	}

	return &out, nil
}
