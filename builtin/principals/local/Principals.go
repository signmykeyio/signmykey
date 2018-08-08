package local

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Principals struct represents map of principals by user.
type Principals struct {
	UserMap map[string]string
}

// Init method is used to ingest config of Principals
func (p *Principals) Init(config *viper.Viper) error {
	p.UserMap = config

	return nil
}

// Get method is used to get the list of principals associated to a specific user.
func (p Principals) Get(user string) ([]string, error) {

	principals := []string{}

	list, ok := p.UserMap[user]
	if !ok {
		return principals, fmt.Errorf("No principals found for %s", user)
	}

	for _, str := range strings.Split(list, ",") {
		trimmed := strings.Trim(str, " ")
		if len(trimmed) > 0 {
			principals = append(principals, trimmed)
		}
	}

	if len(principals) == 0 {
		return principals, fmt.Errorf("No more principals after trim for %s", user)
	}

	return principals, nil
}
