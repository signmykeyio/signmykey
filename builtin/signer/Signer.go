package signer

import (
	"fmt"

	"github.com/spf13/viper"
)

// Signer is the interface that wrap the SMK SSH Signing operation.
type Signer interface {
	Init(config *viper.Viper) error
	Sign(req CertReq) (string, error)
	ReadCA() (string, error)
}

// CertReq represents certificate request
type CertReq struct {
	Key        string
	ID         string
	Principals []string
}

// CheckCertReqFields checks if all CertReq fields are set
func (c CertReq) CheckCertReqFields() error {
	if len(c.Principals) == 0 {
		return fmt.Errorf("Empty list of principals")
	}
	if len(c.ID) == 0 {
		return fmt.Errorf("Empty ID string")
	}
	if len(c.Key) == 0 {
		return fmt.Errorf("Empty key string")
	}

	return nil
}
