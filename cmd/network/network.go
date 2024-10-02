package network

import (
	"fmt"

	"github.com/EliasManj/go-wallet/wallet"
	"github.com/spf13/cobra"
)

func printNetwork(network wallet.Network) {
	if network.Selected {
		fmt.Println("Label:", network.Label, "(Selected)")
	} else {
		fmt.Println("Label:", network.Label)
	}
	fmt.Println("Chain ID: ", network.ChainId)
	fmt.Println("RPC URL: ", network.RpcUrl)
	fmt.Println("Symbol: ", network.Symbol)
	fmt.Println("------------------------------------------------------------------------------------------")
}

var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Network is a palette that contains network based commands",
	Long:  `Manage and configure network settings, you can add any EVM-compatible network to your wallet`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
