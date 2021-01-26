package signer

import (
	"context"

	"github.com/spf13/viper"
)

// Signer is the interface that wrap the SMK SSH Signing operation.
type Signer interface {
	Init(config *viper.Viper) error
	Sign(ctx context.Context, payload []byte, id string, principals []string) (cert string, err error)
	ReadCA(ctx context.Context) (cert string, err error)
}

// CertReq represents certificate request
type CertReq struct {
	Key        string
	ID         string
	Principals []string
}
