package local

import (
	"errors"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Authenticator struct represents local Authenticator options
type Authenticator struct {
	UserMap *viper.Viper
}

// Init method is used to ingest config of Authenticator
func (a *Authenticator) Init(config *viper.Viper) error {

	if !config.IsSet("users") {
		return errors.New("Missing config entry \"users\" for Authenticator")
	}

	a.UserMap = config.Sub("users")

	return nil
}

// Login method is used to check if a couple of user/password is valid in local config
func (a Authenticator) Login(user, password string) (valid bool, swapuser string, err error) {
	if len(user) == 0 {
		return false, "", errors.New("empty username")
	}
	if len(password) == 0 {
		return false, "", errors.New("empty password")
	}

	hashedPass := a.UserMap.GetString(user)
	if len(hashedPass) == 0 {
		return false, "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		return false, "", errors.New("bad password")
	}

	return true, "", nil
}
