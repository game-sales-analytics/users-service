package passhash

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func Verify(raw, hashed string) (bool, error) {
	return comparePasswordAndHash(raw, hashed)
}

func comparePasswordAndHash(raw, hashed string) (bool, error) {
	p, salt, hash, err := decodeHash(hashed)
	if nil != err {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(raw), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *argon2HashParams, salt, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if nil != err {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &argon2HashParams{}
	if _, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism); nil != err {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if nil != err {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if nil != err {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
