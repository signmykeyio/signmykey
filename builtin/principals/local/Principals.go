package local

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	princsPkg "github.com/signmykeyio/signmykey/builtin/principals"
)

// Principals struct represents map of principals by user.
type Principals struct {
	UserMap *viper.Viper
}

type localPrincipals struct {
	User string `json:"user" binding:"required"`
}

// Init method is used to ingest config of Principals
func (p *Principals) Init(config *viper.Viper) error {
	if !config.IsSet("users") {
		return errors.New("Missing config entry \"users\" for Principals")
	}

	p.UserMap = config.Sub("users")

	return nil
}

// Get method is used to get the list of principals associated to a specific user.
func (p Principals) Get(ctx context.Context, payload []byte) (context.Context, []string, error) {

	var local localPrincipals
	err := json.Unmarshal(payload, &local)
	if err != nil {
		log.Errorf("json unmarshaling failed: %s", err)
		return ctx, []string{}, fmt.Errorf("JSON unmarshaling failed: %w", err)
	}

	if !p.UserMap.IsSet(local.User) {
		return ctx, []string{}, princsPkg.NewNotFoundError("local", "No principals found")
	}

	principals := []string{}
	for _, str := range strings.Split(p.UserMap.GetString(local.User), ",") {
		trimmed := strings.Trim(str, " ")
		if len(trimmed) > 0 {
			principals = append(principals, trimmed)
		}
	}

	if len(principals) == 0 {
		return ctx, principals, princsPkg.NewNotFoundError("local", "No more principals after trim")
	}

	return ctx, principals, nil
}
