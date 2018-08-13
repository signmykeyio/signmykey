package cmd

import (
	"bytes"

	"github.com/signmykeyio/signmykey/api"
	localAuth "github.com/signmykeyio/signmykey/builtin/authenticator/local"
	localPrinc "github.com/signmykeyio/signmykey/builtin/principals/local"
	localSign "github.com/signmykeyio/signmykey/builtin/signer/local"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverDevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start signmykey server in DEV mode",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Log level
		logrus.SetLevel(logrus.InfoLevel)
		logrus.Info("start signmykey server in DEV mode")

		// Authenticator init
		auth := &localAuth.Authenticator{}
		authConfig := viper.New()
		authConfig.SetConfigType("yaml")
		authConfig.ReadConfig(bytes.NewBuffer([]byte("admin: $2a$10$1RHgVN4p0QRv9b4Xb2hdDe8qkokyu7vIG7Cx7sDVtKdAaG52vWbuW")))
		auth.Init(authConfig)

		// Principals init
		princs := &localPrinc.Principals{}

		// Signer init
		signer := &localSign.Signer{}

		config := api.Config{
			Auth:   auth,
			Princs: princs,
			Signer: signer,

			Addr:       "127.0.0.1:9600",
			TLSDisable: true,
		}

		err := api.Serve(config)

		return err
	},
}

func init() {
	serverCmd.AddCommand(serverDevCmd)
}
