package auth

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func generateToken(userID, secret string) (string, error) {
	tokenID, err := uuid.NewV4()
	if nil != err {
		return "", err
	}

	if nil != err {
		return "", err
	}

	token := jwt.New()
	if err := token.Set(jwt.IssuerKey, "https://github.com/game-sales-analytics/users-service"); nil != err {
		return "", errors.New("unable to generate auth token")
	}
	if err := token.Set(jwt.AudienceKey, []string{"users"}); nil != err {
		return "", errors.New("unable to generate auth token")
	}
	if err := token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*24*7)); nil != err {
		return "", errors.New("unable to generate auth token")
	}
	if err := token.Set(jwt.IssuedAtKey, time.Now()); nil != err {
		return "", errors.New("unable to generate auth token")
	}
	if err := token.Set(jwt.JwtIDKey, tokenID); nil != err {
		return "", errors.New("unable to generate auth token")
	}
	if err := token.Set(jwt.NotBeforeKey, time.Now()); nil != err {
		return "", errors.New("unable to generate auth token")
	}
	if err := token.Set(jwt.SubjectKey, userID); nil != err {
		return "", errors.New("unable to generate auth token")
	}

	serialized, err := jwt.Sign(token, jwa.HS512, []byte(secret))

	return string(serialized), err
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
