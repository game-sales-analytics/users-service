package id

import (
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
)

func GenerateUserLoginID() (string, error) {
	return xid.New().String(), nil
}

func GenerateUserID() (string, error) {
	uniqueID, err := ksuid.NewRandom()
	if nil != err {
		return "", nil
	}

	return uniqueID.String(), nil
}
