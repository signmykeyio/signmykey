package local

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"github.com/signmykeyio/signmykey/util"
)

// Authenticator struct represents local Authenticator options
type Authenticator struct {
	UserMap *viper.Viper
}

type localLogin struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	Otp string `json:"otp"`
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
		return ctx, false, "", fmt.Errorf("JSON unmarshaling failed: %w", err)
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

	passAndOtp := strings.Split(hashedPass, ",")
	if len(passAndOtp[1]) !=0 && len(login.Otp) == 0  {
		return ctx, false, "", errors.New("otp required but not provided")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passAndOtp[0]), []byte(login.Password))
	if err != nil {
		return ctx, false, "", errors.New("bad password")
	}

	if len(passAndOtp[1]) !=0 {
		seed := util.DecryptSeed(passAndOtp[1], []byte(login.Password))
		timeval := time.Now().Unix() / 30
		generated := util.GenerateOTPCode(seed, timeval)
		if login.Otp != generated {
			// try again as industry standard is a 30 sec tolerance window
			timeval = timeval - 30
			if login.Otp != util.GenerateOTPCode(seed, timeval) {
				return ctx, false, "", errors.New("otp does not match")
			}
		}
	}

	return ctx, true, fmt.Sprintf("local-%s", login.User), nil
}
