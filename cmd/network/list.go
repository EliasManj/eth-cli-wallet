package network

import (
	"fmt"
	"log"

	"github.com/EliasManj/go-wallet/utils"
	"github.com/EliasManj/go-wallet/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getNetworks() {
	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	networks, err := wallet.ListNetworks(db) // Call listNetworks function
	if err != nil {
		log.Fatalf("Failed to list networks: %v", err)
	}
	fmt.Println("Networks:")
	fmt.Println("")
	for _, network := range networks {
		printNetwork(network)
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "This list all the networks available for the wallet",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		getNetworks()
	},
}

func init() {
	NetworkCmd.AddCommand(listCmd)
}
