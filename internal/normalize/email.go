package normalize

import (
	normalizer "github.com/dimuska139/go-email-normalizer"
)

func Email(email string) (string, error) {
	return normalizer.NewNormalizer().Normalize(email), nil
}
