package account

import (
	"fmt"
	"log"

	"github.com/EliasManj/go-wallet/utils"
	"github.com/EliasManj/go-wallet/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func showBalance() {
	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	account, err := wallet.GetSelectedAccount(db)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}

	network, err := wallet.GetSelectedNetwork(db)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	balance, err := wallet.GetBalance(account.Publicy, network)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	fmt.Printf("Balance: %s %s", balance.String(), network.Symbol)
	fmt.Println()

}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "This command displays the balance of the selected account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		showBalance()
	},
}

func init() {
	AccountCmd.AddCommand(balanceCmd)
}
