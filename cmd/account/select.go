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
	label_select string
)

func selectAccount() {
	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	account, err := wallet.GetAccount(db, label_select)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}

	err = wallet.SelectAccount(db, account.Label)
	if err != nil {
		log.Fatalf("Failed to select account: %v", err)
	}

	fmt.Printf("Account %s selected\n", account.Label)
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "This selects the specified account to be the active account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		selectAccount()
	},
}

func init() {
	AccountCmd.AddCommand(selectCmd)
	selectCmd.Flags().StringVarP(&label_select, "label", "l", "", "Label for the network to be identified with")
	selectCmd.MarkFlagRequired("label")
}
