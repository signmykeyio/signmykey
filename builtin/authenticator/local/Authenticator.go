package local

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Authenticator struct represents local Authenticator options
type Authenticator struct {
	UserMap *viper.Viper
}

type localLogin struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
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
func (a Authenticator) Login(ctx context.Context, payload []byte) (resultCtx context.Context, valid bool, id string, err error) {

	var login localLogin
	err = json.Unmarshal(payload, &login)
	if err != nil {
		log.Errorf("json unmarshaling failed: %s", err)
		return ctx, false, "", fmt.Errorf("JSON unmarshaling failed: %s", err)
	}

	if len(login.User) == 0 {
		return ctx, false, "", errors.New("empty username")
	}
	if len(login.Password) == 0 {
		return ctx, false, "", errors.New("empty password")
	}

	hashedPass := a.UserMap.GetString(login.User)
	if len(hashedPass) == 0 {
		return ctx, false, "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(login.Password))
	if err != nil {
		return ctx, false, "", errors.New("bad password")
	}

	return ctx, true, fmt.Sprintf("local-%s", login.User), nil
}
