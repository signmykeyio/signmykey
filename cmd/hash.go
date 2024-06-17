package cmd

import (
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"github.com/signmykeyio/signmykey/util"
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

		fmt.Printf("\nDo you want to use One Time Codes e.g. Google authenticator? (Y/N) ")
		var useOtp string
		_, err = fmt.Scanln(&useOtp)
		if err != nil {
			return err
		}

		if useOtp == "Y" || useOtp == "y" {
			seed := util.GenerateSeed()
			encryptedSeed := util.EncryptSeed(seed, password)
			str := util.ProvisionURI(seed)

			fmt.Printf("\nScan this with your OTP application\n")
			qrterminal.GenerateHalfBlock(str, qrterminal.L, os.Stdout)
			fmt.Printf("\n...or create a new OTP secret manually if you cannot scan QR codes")
			fmt.Printf("\nOTP Secret: %s\n", seed)
			fmt.Printf("\nHashed password: %s\n", string(hash)+","+encryptedSeed)

		} else {
			fmt.Printf("\nHashed password: %s\n", string(hash))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
}
