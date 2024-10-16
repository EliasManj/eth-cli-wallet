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
	label_remove string
)

func removeAccount() {
	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	account, err := wallet.GetAccount(db, label_remove)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}

	err = wallet.RemoveAccount(db, account.Label)
	if err != nil {
		log.Fatalf("Failed to delete account: %v", err)
	}

	fmt.Printf("Account %s deleted\n", account.Label)
}

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "This command removes an account with the given label",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		removeAccount()
	},
}

func init() {
	AccountCmd.AddCommand(removeCmd)
	removeCmd.Flags().StringVarP(&label_remove, "label", "l", "", "Label for the account to be deleted")
	removeCmd.MarkFlagRequired("label")
}
