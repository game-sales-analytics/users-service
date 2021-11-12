package normalize

import (
	normalizer "github.com/dimuska139/go-email-normalizer"
)

func NormalizeEmail(email string) (string, error) {
	return normalizer.NewNormalizer().Normalize(email), nil
}
