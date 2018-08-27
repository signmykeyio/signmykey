package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/signmykeyio/signmykey/api"
	"github.com/signmykeyio/signmykey/builtin/authenticator"
	ldapAuth "github.com/signmykeyio/signmykey/builtin/authenticator/ldap"
	localAuth "github.com/signmykeyio/signmykey/builtin/authenticator/local"
	"github.com/signmykeyio/signmykey/builtin/principals"
	ldapPrinc "github.com/signmykeyio/signmykey/builtin/principals/ldap"
	localPrinc "github.com/signmykeyio/signmykey/builtin/principals/local"
	"github.com/signmykeyio/signmykey/builtin/signer"
	localSign "github.com/signmykeyio/signmykey/builtin/signer/local"
	vaultSign "github.com/signmykeyio/signmykey/builtin/signer/vault"
	"github.com/signmykeyio/signmykey/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCfgFile string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start signmykey server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// load config
		if err := initConfig(serverCfgFile); err != nil {
			return err
		}

		logger := logrus.New()
		logger.Formatter = &logging.TextFormatter{}

		logger.Info("start signmykey server")

		// Log level
		logLevelConfig := map[string]logrus.Level{
			"debug": logrus.DebugLevel,
			"info":  logrus.InfoLevel,
			"warn":  logrus.WarnLevel,
			"fatal": logrus.FatalLevel,
			"panic": logrus.PanicLevel,
		}
		viper.SetDefault("logLevel", "info")
		logLevel, ok := logLevelConfig[strings.ToLower(viper.GetString("logLevel"))]
		if !ok {
			logrus.Fatalf("invalid logLevel %s (debug,info,warn,fatal,panic)", viper.GetString("logLevel"))
		}
		logger.SetLevel(logLevel)

		// Authenticator init
		authTypeConfig := viper.GetString("authenticatorType")
		if authTypeConfig == "" {
			logrus.Fatal("authenticator type not defined in config")
		}
		authType := map[string]authenticator.Authenticator{
			"ldap":  &ldapAuth.Authenticator{},
			"local": &localAuth.Authenticator{},
		}
		auth, ok := authType[authTypeConfig]
		if !ok {
			return fmt.Errorf("unknown authenticator type %s", authTypeConfig)
		}
		err := auth.Init(viper.Sub("authenticatorOpts"), logger)
		if err != nil {
			return err
		}

		// Principals init
		princsTypeConfig := viper.GetString("principalsType")
		if princsTypeConfig == "" {
			return errors.New("principals type not defined in config")
		}
		princsType := map[string]principals.Principals{
			"local": &localPrinc.Principals{},
			"ldap":  &ldapPrinc.Principals{},
		}
		princs, ok := princsType[princsTypeConfig]
		if !ok {
			return fmt.Errorf("unknown principals type %s", princsTypeConfig)
		}
		err = princs.Init(viper.Sub("principalsOpts"))
		if err != nil {
			return err
		}

		// Signer init
		signerTypeConfig := viper.GetString("signerType")
		if signerTypeConfig == "" {
			return errors.New("signer type not defined in config")
		}
		signerType := map[string]signer.Signer{
			"vault": &vaultSign.Signer{},
			"local": &localSign.Signer{},
		}
		signer, ok := signerType[signerTypeConfig]
		if !ok {
			return fmt.Errorf("unknown signer type %s", signerTypeConfig)
		}
		err = signer.Init(viper.Sub("SignerOpts"))
		if err != nil {
			return err
		}

		viper.SetDefault("address", "0.0.0.0:9600")
		viper.SetDefault("tlsDisable", false)

		if !viper.GetBool("tlsDisable") {
			if viper.GetString("tlsCert") == "" || viper.GetString("tlsKey") == "" {
				return fmt.Errorf("tlsCert and tlsKey must be defined if tlsDisable is False")
			}
		}

		config := api.Config{
			Auth:   auth,
			Princs: princs,
			Signer: signer,

			Addr:       viper.GetString("address"),
			TLSDisable: viper.GetBool("tlsDisable"),
			TLSCert:    viper.GetString("tlsCert"),
			TLSKey:     viper.GetString("tlsKey"),
		}

		err = api.Serve(config)

		return err
	},
}

func init() {
	serverCmd.Flags().StringVarP(
		&serverCfgFile, "cfg", "c", "/etc/signmykey/server.yml", "config file")

	rootCmd.AddCommand(serverCmd)
}
