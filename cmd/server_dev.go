package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/fatih/color"
	"github.com/signmykeyio/signmykey/api"
	localAuth "github.com/signmykeyio/signmykey/builtin/authenticator/local"
	localPrinc "github.com/signmykeyio/signmykey/builtin/principals/local"
	localSign "github.com/signmykeyio/signmykey/builtin/signer/local"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
)

var devUser string

var serverDevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start signmykey server in DEV mode",
	Run: func(cmd *cobra.Command, args []string) {

		// Logging
		logger := logrus.New()
		logger.SetLevel(logrus.InfoLevel)
		logger.WithField("ctx", "server").Warn("Starting signmykey server in DEV mode")

		// Authenticator init
		password, hash, err := generateAndHashPassword()
		if err != nil {
			logger.WithField("ctx", "server").WithError(err).Error("Generating new password and hash")
			return
		}
		auth := &localAuth.Authenticator{}
		authConfig := viper.New()
		authConfig.SetConfigType("yaml")
		err = authConfig.ReadConfig(bytes.NewBuffer([]byte(fmt.Sprintf("%s: %s", devUser, hash))))
		if err != nil {
			logger.WithField("ctx", "server").WithError(err).Error("Loading local authenticator config")
			return
		}
		auth.UserMap = authConfig

		// Principals init
		princs := &localPrinc.Principals{}
		princsConfig := viper.New()
		princsConfig.SetConfigType("yaml")
		err = princsConfig.ReadConfig(bytes.NewBuffer([]byte(fmt.Sprintf(`
users:
  %s: %s`, devUser, devUser))))
		if err != nil {
			logger.WithField("ctx", "server").WithError(err).Error("Loading local principals config")
			return
		}
		err = princs.Init(princsConfig)
		if err != nil {
			logger.WithField("ctx", "server").WithError(err).Error("initializing local principals")
			return
		}

		// Signer init
		caSigner, caPub, err := generateCA()
		if err != nil {
			logger.WithField("ctx", "server").WithError(err).Error("Generating temp SSH Certificate Authority")
			return
		}
		signer := &localSign.Signer{
			CACert: caPub,
			CAKey:  caSigner,
			TTL:    600,
			Extensions: map[string]string{
				"permit-X11-forwarding":   "",
				"permit-agent-forwarding": "",
				"permit-port-forwarding":  "",
				"permit-pty":              "",
				"permit-user-rc":          "",
			},
		}

		config := api.Config{
			Auth:   auth,
			Princs: princs,
			Signer: signer,

			Logger: logger,

			Addr:       "127.0.0.1:9600",
			TLSDisable: true,
		}

		displayHowto(password, ssh.MarshalAuthorizedKey(caPub))

		api.Serve(config)
	},
}

func generateCA() (ssh.Signer, ssh.PublicKey, error) {
	privateSeed, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating CA private key: %w", err)
	}

	privateBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateSeed),
	}

	signer, err := ssh.ParsePrivateKey(pem.EncodeToMemory(&privateBlock))
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing CA private key: %w", err)
	}

	public, err := ssh.NewPublicKey(&privateSeed.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("error generation CA public key: %w", err)
	}

	return signer, public, nil
}

func generateAndHashPassword() (string, string, error) {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	passwordBytes := make([]byte, 30)

	for i := range passwordBytes {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", "", fmt.Errorf("error getting random number for password generation: %w", err)
		}
		passwordBytes[i] = letterBytes[random.Int64()]
	}

	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("error hashing generated password: %w", err)
	}

	return string(passwordBytes), string(hash), nil
}

func displayHowto(password string, ca []byte) {

	color.Red(`
WARNING! Dev mode is enabled! In this mode, Signmykey runs
WARNING! entirely in-memory so data are not persisted!
WARNING! Dev mode should NOT be used in production installations!
`)

	color.Blue("\n### Server side\n")
	color.Yellow(`
An ephemeral certificate authority is created for this instance and will die with it.
To deploy this CA on destination servers, you can launch this command:

	$ echo "%s" > /etc/ssh/ca.pub

You then have to add this line to "/etc/ssh/sshd_config" and restart OpenSSH server:

	TrustedUserCAKeys /etc/ssh/ca.pub
`, string(ca)[0:(len(ca)-1)])

	color.Blue("\n### Client side\n")
	color.Yellow(`
A temporary user is created with this parameters:

	user: %s
	password: %s
	principals: %s

You can sign your key with this command:

	$ signmykey -a http://127.0.0.1:9600/ -u %s


`, devUser, password, devUser, devUser)
}

func init() {
	serverDevCmd.Flags().StringVarP(
		&devUser, "user", "u", "admin", "ephemeral user to use with Dev mode")

	serverCmd.AddCommand(serverDevCmd)
}
