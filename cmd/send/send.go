package send

import (
	"fmt"
	"log"
	"math/big"

	"github.com/EliasManj/go-wallet/utils"
	"github.com/EliasManj/go-wallet/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	amount_send string
	to_send     string
)

func sendFunction() {

	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	account, err := wallet.GetSelectedAccount(db)
	if err != nil {
		fmt.Printf("Failed to get account: %v\n", err)
		return
	}

	network, err := wallet.GetSelectedNetwork(db)
	if err != nil {
		fmt.Printf("Failed to get network: %v\n", err)
		return
	}

	amount := new(big.Int)
	amount, ok := amount.SetString(amount_send, 10)
	if !ok {
		fmt.Printf("Failed to convert amount to big.Int\n")
		return
	}

	fmt.Printf("Sending from account %s on network %s\n", account.Label, network.Label)

	tx, err := wallet.SendETH(account.Privatey, to_send, amount, network)
	if err != nil {
		fmt.Printf("Failed to send ETH: %v\n", err)
		return
	}

	printTx(tx)

}

func printTx(tx wallet.Transaction) {
	fmt.Printf("Transaction hash: %s\n", tx.Hash)
	fmt.Printf("From: %s\n", tx.From)
	fmt.Printf("To: %s\n", tx.To)
	fmt.Printf("Amount: %s\n", tx.Amount)
	fmt.Printf("Gas: %s\n", tx.GasUsed)
	fmt.Printf("Gas price: %s\n", tx.GasPrice)
}

var SendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send ether or tokens from the selected address and network to the specified account",
	Long:  `Send ether or tokens from the selected address and network to the specified account`,
	Run: func(cmd *cobra.Command, args []string) {
		sendFunction()
	},
}

func init() {
	SendCmd.Flags().StringVarP(&to_send, "to", "t", "", "Label for the network to be identified with")
	SendCmd.MarkFlagRequired("to")
	SendCmd.Flags().StringVarP(&amount_send, "amt", "a", "", "Label for the network to be identified with")
	SendCmd.MarkFlagRequired("amt")
}
