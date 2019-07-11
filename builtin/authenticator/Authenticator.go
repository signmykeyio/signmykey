package authenticator

import (
	"context"

	"github.com/spf13/viper"
)

// Authenticator is the interface that wrap the SMK Authentication logic.
type Authenticator interface {
	Init(config *viper.Viper) error
	Login(ctx context.Context, payload []byte) (resultCtx context.Context, valid bool, id string, err error)
}
