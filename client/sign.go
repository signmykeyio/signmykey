package client

import (
	"errors"

	"github.com/dghubble/sling"
)

type signLDAPRequest struct {
	User      string `json:"user"`
	Password  string `json:"password"`
	PublicKey string `json:"public_key"`
}

type signLDAPResponse struct {
	Certificate string `json:"certificate"`
}

type signLDAPError struct {
	Error string `json:"error"`
}

// Sign is used to sign an SSH key with user/password combination.
func Sign(addr, user, password, pubKey string) (certificate string, err error) {

	body := &signLDAPRequest{
		User:      user,
		Password:  password,
		PublicKey: pubKey,
	}

	signResponse := &signLDAPResponse{}
	signError := &signLDAPError{}

	res, err := sling.New().Post(addr).Path("v1/sign").BodyJSON(body).Receive(signResponse, signError)
	if err != nil {
		return certificate, err
	}

	if res.StatusCode != 200 {
		err = errors.New(signError.Error)
		return certificate, err
	}

	certificate = signResponse.Certificate
	return certificate, err
}
