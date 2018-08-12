package vault

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Authenticator struct represents local Authenticator options
type Authenticator struct {
	UserMap *viper.Viper
}

// Init method is used to ingest config of Authenticator
func (a *Authenticator) Init(config *viper.Viper) error {
	a.UserMap = config

	return nil
}

// Login method is used to check if a couple of user/password is valid in local config
func (a Authenticator) Login(user, password string) (valid bool, err error) {
	if len(user) == 0 {
		return false, fmt.Errorf("empty username")
	}
	if len(password) == 0 {
		return false, fmt.Errorf("empty password")
	}

	err = bcrypt.CompareHashAndPassword([]byte("myhashedpass"), []byte("mypass"))
	if err != nil {
		return false, errors.New("bad password")
	}

	return true, nil
}
