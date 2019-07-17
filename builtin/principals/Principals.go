package principals

import (
	"context"

	"github.com/spf13/viper"
)

// Principals is the interface that wrap the get of SMK Principals.
type Principals interface {
	Init(config *viper.Viper) error
	Get(ctx context.Context, payload []byte) (context.Context, []string, error)
}
