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

// NotFoundError it's principals provider error when no principals is found
type NotFoundError struct {
	provider string
	msg      string
}

// NewNotFoundError creates new NotFoundError error with given privder name and message
func NewNotFoundError(provider, msg string) *NotFoundError {
	return &NotFoundError{provider: provider, msg: msg}
}

func (e *NotFoundError) Error() string {
	return e.provider + ": " + e.msg
}
