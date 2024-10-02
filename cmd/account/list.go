package account

import (
	"fmt"
	"log"

	"github.com/EliasManj/go-wallet/utils"
	"github.com/EliasManj/go-wallet/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func listAccounts() {
	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	accounts, err := wallet.ListAccounts(db)
	if err != nil {
		log.Fatalf("Failed to list accounts: %v", err)
	}
	fmt.Println("Accounts:")
	fmt.Println("")
	for _, account := range accounts {
		printAccount(account)
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "This list all the networks available for the wallet",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		listAccounts()
	},
}

func init() {
	AccountCmd.AddCommand(listCmd)
}
