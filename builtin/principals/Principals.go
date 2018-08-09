package principals

import "github.com/spf13/viper"

// Principals is the interface that wrap the get of SMK Principals.
type Principals interface {
	Init(config *viper.Viper) error
	Get(user string) ([]string, error)
}
