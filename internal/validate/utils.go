package validate

import (
	"net"
	"net/mail"
)

func isEmailValid(email string) (bool, error) {
	_, err := mail.ParseAddress(email)

	return err == nil, nil
}

func isIPValid(ip string) (bool, error) {
	return net.ParseIP(ip) != nil, nil
}
