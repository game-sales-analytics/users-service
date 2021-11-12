package passhash

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func defaultArgon2HashParams() argon2HashParams {
	return argon2HashParams{
		memory:      128 * 1024,
		iterations:  128,
		parallelism: 8,
		saltLength:  2048,
		keyLength:   4096,
	}
}

type argon2HashParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func HashPassword(raw string) (string, error) {
	params := defaultArgon2HashParams()
	return generateFromPassword(raw, params)
}

func generateFromPassword(password string, params argon2HashParams) (string, error) {
	salt, err := generateRandomBytes(params.saltLength)
	if nil != err {
		return "", ErrGenerateRandom
	}

	hash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.memory,
		params.iterations,
		params.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); nil != err {
		return nil, err
	}

	return b, nil
}
