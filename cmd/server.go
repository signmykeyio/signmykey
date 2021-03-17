package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/signmykeyio/signmykey/api"
	"github.com/signmykeyio/signmykey/builtin/authenticator"
	ldapAuth "github.com/signmykeyio/signmykey/builtin/authenticator/ldap"
	localAuth "github.com/signmykeyio/signmykey/builtin/authenticator/local"
	oidcropcAuth "github.com/signmykeyio/signmykey/builtin/authenticator/oidcropc"
	"github.com/signmykeyio/signmykey/builtin/principals"
	ldapPrinc "github.com/signmykeyio/signmykey/builtin/principals/ldap"
	localPrinc "github.com/signmykeyio/signmykey/builtin/principals/local"
	oidcropcPrinc "github.com/signmykeyio/signmykey/builtin/principals/oidcropc"
	userPrinc "github.com/signmykeyio/signmykey/builtin/principals/user"
	"github.com/signmykeyio/signmykey/builtin/signer"
	localSign "github.com/signmykeyio/signmykey/builtin/signer/local"
	vaultSign "github.com/signmykeyio/signmykey/builtin/signer/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverCfgFile   string
	serverLogFormat string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start signmykey server",
	Run: func(cmd *cobra.Command, args []string) {

		// Create logger
		logFormatter := map[string]logrus.Formatter{
			"json": &logrus.JSONFormatter{DisableTimestamp: true},
			"text": &logrus.TextFormatter{},
		}
		if serverLogFormat != "json" && serverLogFormat != "text" {
			fmt.Println("Flag --log-format value must be \"json\" or \"text\"")
			os.Exit(1)
		}
		logger := logrus.New()
		logger.Formatter = logFormatter[serverLogFormat]

		// load config
		if err := initConfig(serverCfgFile); err != nil {
			logger.WithField("ctx", "server").WithError(fmt.Errorf("failed to load config! %s", err)).Error("Loading config")
			return
		}

		logger.WithField("ctx", "server").Info("Starting signmykey server")

		// Log level
		viper.SetDefault("logLevel", "info")
		logLevel, err := logrus.ParseLevel(viper.GetString("logLevel"))
		if err != nil {
			logger.WithField("ctx", "server").WithError(err).Error("Setting logLevel")
			return
		}
		logger.SetLevel(logLevel)

		// Authenticator init
		authTypeConfig := viper.GetString("authenticatorType")
		if authTypeConfig == "" {
			logger.WithField("ctx", "server").WithError(errors.New("authenticator type not defined in config")).Error("Setting Authenticator type")
			return
		}
		authType := map[string]authenticator.Authenticator{
			"local":    &localAuth.Authenticator{},
			"ldap":     &ldapAuth.Authenticator{},
			"oidcropc": &oidcropcAuth.Authenticator{},
		}
		auth, ok := authType[authTypeConfig]
		if !ok {
			logger.WithField("ctx", "server").WithError(fmt.Errorf("unknown authenticator type %s", authTypeConfig)).Error("Setting Authenticator type")
			return
		}
		err = auth.Init(viper.Sub("authenticatorOpts"))
		if err != nil {
			logger.WithField("ctx", "server").WithError(err).Error("Setting Authenticator options")
			return
		}

		// Principals init
		princsType := map[string]principals.Principals{
			"local":    &localPrinc.Principals{},
			"ldap":     &ldapPrinc.Principals{},
			"oidcropc": &oidcropcPrinc.Principals{},
			"user":     &userPrinc.Principals{},
		}
		princsProviders := []principals.Principals{}

		if viper.IsSet("principalsProviders") {
			for princsTypeConfig := range viper.GetStringMap("principalsProviders") {
				princsOptsSection := "principalsProviders." + princsTypeConfig

				logger.WithField("ctx", "server").Infof("Configure %v principals provider", princsTypeConfig)
				princs, ok := princsType[princsTypeConfig]
				if !ok {
					logger.WithField("ctx", "server").WithError(fmt.Errorf("unknown principals type %s", princsTypeConfig)).Error("Setting Principals type")
					return
				}
				err = princs.Init(viper.Sub(princsOptsSection))
				if err != nil {
					logger.WithField("ctx", "server").WithError(err).Error("Setting Principals options")
					return
				}

				princsProviders = append(princsProviders, princs)
			}
		} else {
			princsTypeConfig := viper.GetString("principalsType")
			if princsTypeConfig == "" {
				logger.WithField("ctx", "server").WithError(errors.New("principals type not defined in config")).Error("Setting Principals type")
				return
			}

			logger.WithField("ctx", "server").Infof("Configure %v principals provider", princsTypeConfig)
			princs, ok := princsType[princsTypeConfig]
			if !ok {
				logger.WithField("ctx", "server").WithError(fmt.Errorf("unknown principals type %s", princsTypeConfig)).Error("Setting Principals type")
				return
			}
			err = princs.Init(viper.Sub("principalsOpts"))
			if err != nil {
				logger.WithField("ctx", "server").WithError(err).Error("Setting Principals options")
				return
			}

			princsProviders = append(princsProviders, princs)
		}

		if len(princsProviders) == 0 {
			logger.WithField("ctx", "server").Error("principals providers list is not configured")
			return
		}

		// Signer init
		signerTypeConfig := viper.GetString("signerType")
		if signerTypeConfig == "" {
			logger.WithField("ctx", "server").WithError(errors.New("singer type not defined in config")).Error("Setting Signer type")
			return
		}
		signerType := map[string]signer.Signer{
			"vault": &vaultSign.Signer{},
			"local": &localSign.Signer{},
		}
		signer, ok := signerType[signerTypeConfig]
		if !ok {
			logger.WithField("ctx", "server").WithError(fmt.Errorf("unknown signer type %s", signerTypeConfig)).Error("Setting Signer type")
			return
		}
		err = signer.Init(viper.Sub("SignerOpts"))
		if err != nil {
			logger.WithField("ctx", "server").WithError(err).Error("Setting Signer options")
			return
		}

		viper.SetDefault("address", "0.0.0.0:9600")
		viper.SetDefault("tlsDisable", false)

		if !viper.GetBool("tlsDisable") {
			if viper.GetString("tlsCert") == "" || viper.GetString("tlsKey") == "" {
				logger.WithField("ctx", "server").WithError(errors.New("tlsCert and tlsKey must be defined if tlsDisable is False")).Error("Setting TLS config")
				return
			}
		}

		config := api.Config{
			Auth:   auth,
			Princs: princsProviders,
			Signer: signer,

			Logger: logger,

			Addr:       viper.GetString("address"),
			TLSDisable: viper.GetBool("tlsDisable"),
			TLSCert:    viper.GetString("tlsCert"),
			TLSKey:     viper.GetString("tlsKey"),
		}

		api.Serve(config)
		logger.WithField("ctx", "server").Info("Stopping HTTP server")
	},
}

func init() {
	serverCmd.Flags().StringVarP(&serverCfgFile, "cfg", "c", "/etc/signmykey/server.yml", "config file")
	serverCmd.Flags().StringVarP(&serverLogFormat, "log-format", "l", "json", "logging format (json/text)")

	rootCmd.AddCommand(serverCmd)
}
