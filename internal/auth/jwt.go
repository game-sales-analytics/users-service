package auth

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type GeneratedToken struct {
	ID                 string
	Value              string
	NotBeforeDateTime  time.Time
	ExpirationDateTime time.Time
}

func (a *authsrv) generateToken(userID, secret string) (*GeneratedToken, error) {
	tokenID, err := uuid.NewV4()
	if nil != err {
		a.logger.WithError(err).WithField("err_code", "E_GENERATE_JWT_TOKEN_ID").Error("failed generating jwt token id")
		return nil, err
	}

	token := jwt.New()
	if err := token.Set(jwt.IssuerKey, "https://github.com/game-sales-analytics/users-service"); nil != err {
		a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_ISSUER_CLAIM").Error("failed setting jwt token issuer claim")
		return nil, errors.New("unable to set auth token issuer key")
	}

	if err := token.Set(jwt.AudienceKey, []string{"users"}); nil != err {
		a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_AUDIENCE_CLAIM").Error("failed setting jwt token audience claim")
		return nil, errors.New("unable to set auth token audience key")
	}

	exp := time.Now()
	if err := token.Set(jwt.ExpirationKey, exp.Add(time.Hour*24*7)); nil != err {
		a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_EXPIRATION_CLAIM").Error("failed setting jwt token expiration claim")
		return nil, errors.New("unable to set auth token expiration key")
	}

	if err := token.Set(jwt.IssuedAtKey, time.Now()); nil != err {
		a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_ISSUED_AT_CLAIM").Error("failed setting jwt token issued_at claim")
		return nil, errors.New("unable to set auth token issued_at key")
	}

	if err := token.Set(jwt.JwtIDKey, tokenID.String()); nil != err {
		a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_JWT_ID_CLAIM").Error("failed setting jwt token jwt_id claim")
		return nil, errors.New("unable to set auth token jwt_id key")
	}

	nbf := time.Now()
	if err := token.Set(jwt.NotBeforeKey, nbf); nil != err {
		a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_NOT_BEFORE_CLAIM").Error("failed setting jwt token not_before claim")
		return nil, errors.New("unable to set auth token not_before key")
	}

	if err := token.Set(jwt.SubjectKey, userID); nil != err {
		a.logger.WithError(err).WithField("err_code", "E_SET_JWT_TOKEN_SUBJECT_CLAIM").Error("failed setting jwt token subject claim")
		return nil, errors.New("unable to set auth token subject key")
	}

	serialized, err := jwt.Sign(token, jwa.HS512, []byte(secret))
	if nil != err {
		a.logger.WithError(err).WithField("err_code", "E_SIGN_JWT_TOKEN").Error("failed signing jwt token")
		return nil, errors.New("unable to sign auth token")
	}

	return &GeneratedToken{
		ID:                 tokenID.String(),
		Value:              string(serialized),
		NotBeforeDateTime:  nbf,
		ExpirationDateTime: exp,
	}, err
}

func verifyToken(token, secret string) (*TokenVerificationResult, error) {
	parsedToken, err := jwt.Parse([]byte(token), jwt.WithVerify(jwa.HS512, []byte(secret)))
	if nil != err {
		return nil, err
	}

	out := TokenVerificationResult{
		UserID: parsedToken.Subject(),
	}

	return &out, nil
}
