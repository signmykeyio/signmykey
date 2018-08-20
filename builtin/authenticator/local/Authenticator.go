package local

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Authenticator struct represents local Authenticator options
type Authenticator struct {
	Logger  logrus.FieldLogger
	UserMap *viper.Viper
}

// Init method is used to ingest config of Authenticator
func (a *Authenticator) Init(config *viper.Viper, logger logrus.FieldLogger) error {
	a.UserMap = config
	a.Logger = logger.WithField("app", "auth")

	return nil
}

// Login method is used to check if a couple of user/password is valid in local config
func (a Authenticator) Login(user, password string) (valid bool, err error) {
	if len(user) == 0 {
		a.Logger.Warnf("empty username")
		return false, errors.New("empty username")
	}
	if len(password) == 0 {
		return false, errors.New("empty password")
	}

	hashedPass := a.UserMap.GetString(user)
	if len(hashedPass) == 0 {
		a.Logger.Warnf("user %s not found", user)
		return false, fmt.Errorf("user %s not found", user)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		a.Logger.Warnf("invalid password for user %s", user)
		return false, fmt.Errorf("invalid password for user %s", user)
	}

	return true, nil
}
