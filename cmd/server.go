package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/signmykey/signmykey/api"
	"gitlab.com/signmykey/signmykey/builtin/authenticator"
	ldapAuth "gitlab.com/signmykey/signmykey/builtin/authenticator/ldap"
	vaultAuth "gitlab.com/signmykey/signmykey/builtin/authenticator/vault"
	"gitlab.com/signmykey/signmykey/builtin/principals"
	ldapPrinc "gitlab.com/signmykey/signmykey/builtin/principals/ldap"
	localPrinc "gitlab.com/signmykey/signmykey/builtin/principals/local"
	"gitlab.com/signmykey/signmykey/builtin/signer"
	vaultSign "gitlab.com/signmykey/signmykey/builtin/signer/vault"
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

		logrus.Info("start signmykey-server")

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
		logrus.SetLevel(logLevel)

		// Authenticator init
		authTypeConfig := viper.GetString("authenticatorType")
		if authTypeConfig == "" {
			logrus.Fatal("authenticator type not defined in config")
		}
		authType := map[string]authenticator.Authenticator{
			"ldap":  &ldapAuth.Authenticator{},
			"vault": &vaultAuth.Authenticator{},
		}
		auth, ok := authType[authTypeConfig]
		if !ok {
			return fmt.Errorf("unknown authenticator type %s", authTypeConfig)
		}
		err := auth.Init(viper.GetStringMapString("authenticatorOpts"))
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
		err = princs.Init(viper.GetStringMapString("principalsOpts"))
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
		}
		signer, ok := signerType[signerTypeConfig]
		if !ok {
			return fmt.Errorf("unknown signer type %s", signerTypeConfig)
		}
		err = signer.Init(viper.GetStringMapString("SignerOpts"))
		if err != nil {
			return err
		}

		config := api.Config{
			TTL:    viper.GetString("ttl"),
			Auth:   auth,
			Princs: princs,
			Signer: signer,
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
