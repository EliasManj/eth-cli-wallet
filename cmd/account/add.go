package account

import (
	"fmt"
	"log"

	"github.com/EliasManj/go-wallet/utils"
	"github.com/EliasManj/go-wallet/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	label string
)

func createAccount() {
	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	fmt.Println(label)
	account, err := wallet.CreateNewAccount(db, label)
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
	}
	printAccount(account)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Add a new account to the wallet",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		createAccount()
	},
}

func init() {
	AccountCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&label, "label", "l", "", "Label for the network to be identified with")
	createCmd.MarkFlagRequired("label")
}
