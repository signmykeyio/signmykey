package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash password to use with local authenticator",
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Printf("Password to hash (will be hidden): ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}

		hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		fmt.Printf("\nHashed password: %s\n", string(hash))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
}
