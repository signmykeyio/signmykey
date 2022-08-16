package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/signmykeyio/signmykey/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var clientCfgFile string

var rootCmd = &cobra.Command{
	Use:           "signmykey",
	Short:         "A client-server to sign ssh keys",
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// load config
		if err := initConfig(clientCfgFile); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		pubKeysFiles := []string{}
		foundPubKeysFiles, err := client.FindUserPubKeys(viper.GetStringSlice("key"))
		if err != nil {
			return err
		}

		if viper.GetBool("expired") {
			// filter only expired keys
			for _, pubKeyFile := range foundPubKeysFiles {
				if !client.CertStillValid(pubKeyFile) {
					pubKeysFiles = append(pubKeysFiles, pubKeyFile)
				}
			}
			if len(pubKeysFiles) == 0 {
				return nil
			}
		} else {
			// process all found pubkeys
			pubKeysFiles = foundPubKeysFiles
		}

		username := viper.GetString("user")
		if username == "" {
			user, err := user.Current()
			if err != nil {
				return err
			}
			username = user.Username
		}

		password := viper.GetString("password")
		if password == "" {
			fmt.Printf("Enter signmykey password (will be hidden): ")
			passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			password = string(passwordBytes)
		}

		smkAddr := viper.GetString("addr")
		if !strings.HasSuffix(smkAddr, "/") {
			smkAddr = smkAddr + "/"
		}

		deprecatedAlgos := []string{}
		for _, pubKeyFile := range pubKeysFiles {
			pubKey, err := client.GetUserPubKey(pubKeyFile)
			if err != nil {
				return fmt.Errorf("%v, public key: %v", err, pubKeyFile)
			}

			signedKey, err := client.Sign(smkAddr, username, password, pubKey)
			if err != nil {
				return fmt.Errorf("%v, public key: %v", err, pubKeyFile)
			}

			err = client.WriteUserSignedKey(signedKey, pubKeyFile)
			if err != nil {
				return fmt.Errorf("%v, public key: %v", err, pubKeyFile)
			}

			color.Green("\nYour SSH Key %s is successfully signed !", pubKeyFile)

			principals, before, algo, err := client.CertInfo(signedKey)
			if err != nil {
				return fmt.Errorf("%v, public key: %v", err, pubKeyFile)
			}
			color.HiBlack("\n  - Valid until: %s", time.Unix(int64(before), 0))
			color.HiBlack("  - Principals: %s", strings.Join(principals, ","))
			if client.CertAlgoIsDeprecated(algo) {
				color.HiYellow("  - Algorithm is deprecated: %s", algo)
				deprecatedAlgos = append(deprecatedAlgos, algo)
			}
		}

		if len(deprecatedAlgos) != 0 {
			color.HiYellow(`
At least one of signed certificates has type that is deprecated by openssh.
if you are experiencing connection errors, try to update signmykey signer backend.
Also, you can temporary enable deprecated algorithm by adding to ~/.ssh/config :`)
			for _, algo := range deprecatedAlgos {
				color.HiYellow("\n  PubkeyAcceptedKeyTypes +%v", algo)
			}
		}

		return nil
	},
}

// Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red(fmt.Sprintf("Error: %s", err))
		os.Exit(1)
	}
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)

	rootCmd.Flags().StringVarP(&clientCfgFile, "cfg", "c", "~/.signmykey.yml", "config file")

	rootCmd.Flags().StringP("addr", "a", "http://127.0.0.1:9600/", "SMK server address")
	if err := viper.BindPFlag("addr", rootCmd.Flags().Lookup("addr")); err != nil {
		color.Red(fmt.Sprintf("%s", err))
		os.Exit(1)
	}

	rootCmd.Flags().StringP("user", "u", "", "User used to login instead of current")
	if err := viper.BindPFlag("user", rootCmd.Flags().Lookup("user")); err != nil {
		color.Red(fmt.Sprintf("%s", err))
		os.Exit(1)
	}

	rootCmd.Flags().StringP("password", "p", "", "Password used to login")
	if err := viper.BindPFlag("password", rootCmd.Flags().Lookup("password")); err != nil {
		color.Red(fmt.Sprintf("%s", err))
		os.Exit(1)
	}

	rootCmd.Flags().StringSliceP("key", "k", client.DefaultSSHKeys, "Path of public key to sign")
	if err := viper.BindPFlag("key", rootCmd.Flags().Lookup("key")); err != nil {
		color.Red(fmt.Sprintf("%s", err))
		os.Exit(1)
	}

	rootCmd.Flags().BoolP("expired", "e", false, "Sign only if existing key already expired")
	if err := viper.BindPFlag("expired", rootCmd.Flags().Lookup("expired")); err != nil {
		color.Red(fmt.Sprintf("%s", err))
		os.Exit(1)
	}
}

func initConfig(cfgFile string) error {
	viper.SetEnvPrefix("smk")
	viper.AutomaticEnv()

	// expand ~ in file path
	expandedCfgFile, err := homedir.Expand(cfgFile)
	if err != nil {
		return err
	}

	// Use config file defined by flag if exists
	if _, err := os.Stat(expandedCfgFile); err == nil {
		viper.SetConfigFile(expandedCfgFile)
		return viper.ReadInConfig()
	}

	// Use default config file if exists
	if _, err := os.Stat("/etc/signmykey/client.yml"); err == nil {
		viper.SetConfigFile("/etc/signmykey/client.yml")
		return viper.ReadInConfig()
	}

	return nil
}
