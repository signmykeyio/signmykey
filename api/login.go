package api

import "github.com/pkg/errors"

// Login represents credentials expected by SMK server
type Login struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	PubKey   string `json:"public_key" binding:"required"`
}

// Validate fields of login struct
func (l *Login) Validate() error {
	if l.User == "" {
		return errors.New("empty user field")
	}

	if l.Password == "" {
		return errors.New("empty password field")
	}

	if l.PubKey == "" {
		return errors.New("empty public_key field")
	}

	return nil
}
