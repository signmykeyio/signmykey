package user

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Principals struct represents user options.
type Principals struct{}

type userPrincipals struct {
	User string `json:"user" binding:"required"`
}

// Init method is used to ingest config of Principals
func (p *Principals) Init(config *viper.Viper) error {
	return nil
}

// Get method is used to get the list of principals associated to a specific user.
func (p Principals) Get(ctx context.Context, payload []byte) (context.Context, []string, error) {

	var userPrinc userPrincipals
	err := json.Unmarshal(payload, &userPrinc)
	if err != nil {
		log.Errorf("json unmarshaling failed: %s", err)
		return ctx, []string{}, fmt.Errorf("JSON unmarshaling failed: %w", err)
	}

	principals := []string{userPrinc.User}

	return ctx, principals, nil
}
