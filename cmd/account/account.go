package account

import (
	"fmt"

	"github.com/EliasManj/go-wallet/wallet"
	"github.com/spf13/cobra"
)

var AccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Accont is a palette that contains account based commands",
	Long:  `Add and list accounts, you can add any EVM-compatible account to your wallet`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func printAccount(account wallet.Account) {
	if account.Selected {
		fmt.Println("Label:", account.Label, "(Selected)")
	} else {
		fmt.Println("Label:", account.Label)
	}
	fmt.Println("Address: ", account.Publicy)
	fmt.Println("Private Key: ", account.Privatey)
	fmt.Println("Tokens: ", account.Tokens)
	fmt.Println("------------------------------------------------------------------------------------------")
}
