package authenticator

import "github.com/spf13/viper"

// Authenticator is the interface that wrap the SMK Authentication logic.
type Authenticator interface {
	Init(config *viper.Viper) error
	Login(user, password string) (bool, string, error)
}
