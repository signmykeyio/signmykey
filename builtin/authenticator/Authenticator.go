package authenticator

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Authenticator is the interface that wrap the SMK Authentication logic.
type Authenticator interface {
	Init(config *viper.Viper, logger logrus.FieldLogger) error
	Login(user, password string) (bool, error)
}
