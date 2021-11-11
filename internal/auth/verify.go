package auth

import "context"

func (a authsrv) VerifyToken(ctx context.Context, token string) (TokenVerificationResult, error) {
	return verifyToken(token, a.cfg.Key)
}
