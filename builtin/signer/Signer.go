package signer

import "github.com/spf13/viper"

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
