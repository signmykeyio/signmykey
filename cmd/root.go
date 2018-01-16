package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	"gitlab.com/signmykey/signmykey/helper"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "signmykey",
	Short:         "A light command to sign ssh keys with SMK server",
	Long: `A light command to sign ssh keys with SMK server

Config file is in "/etc/signmykey/config.yaml"`,
	Example: `  Sign key in non default path
	
	> signmykey -a https://smkserver -k ~/myrsapubkey.pub`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		addr := viper.GetString("addr")
		if addr[len(addr)-1] != '/' {
			return errors.New("SMK Server address must end with a slash")
		}

		err := helper.UserPubKeyExists(viper.GetString("key"))
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetBool("expired") {
			if helper.CertStillValid(viper.GetString("key")) {
				return nil
			}
		}

		username := viper.GetString("user")
		if username == "" {
			user, err := user.Current()
			if err != nil {
				return err
			}
			username = user.Username
		}

		fmt.Printf("Password (will be hidden): ")
		password, err := gopass.GetPasswd()
		if err != nil {
			return err
		}

		pubKey, err := helper.GetUserPubKey(viper.GetString("key"))
		if err != nil {
			return err
		}

		smkAddr := viper.GetString("addr")
		signedKey, err := helper.Sign(smkAddr, username, string(password), pubKey)
		if err != nil {
			return err
		}

		err = helper.WriteUserSignedKey(signedKey, viper.GetString("key"))
		if err != nil {
			return err
		}

		color.Green("\nYour SSH Key is successfully signed")

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red(fmt.Sprintf("Error: %s", err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringP("addr", "a", "http://127.0.0.1:8080/", "SMK server address")
	if err := viper.BindPFlag("addr", rootCmd.Flags().Lookup("addr")); err != nil {
		color.Red(fmt.Sprintf("%s", err))
		os.Exit(1)
	}

	rootCmd.Flags().StringP("user", "u", "", "User used to login instead of current")
	if err := viper.BindPFlag("user", rootCmd.Flags().Lookup("user")); err != nil {
		color.Red(fmt.Sprintf("%s", err))
		os.Exit(1)
	}

	rootCmd.Flags().StringP("key", "k", "~/.ssh/id_rsa.pub", "Path of public key to sign")
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

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AddConfigPath("/etc/signmykey")
	viper.SetConfigName("config")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig() // nolint: errcheck,gas
}
