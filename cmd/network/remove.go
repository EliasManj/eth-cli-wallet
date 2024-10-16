package network

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

func removeNetwork() {
	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	network, err := wallet.GetNetwork(db, label_remove)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	err = wallet.DeleteNetwork(db, network.Label)
	if err != nil {
		log.Fatalf("Failed to delete network: %v", err)
	}

	fmt.Printf("Network %s deleted\n", network.Label)
}

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "This command removes a network with the given label",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		removeNetwork()
	},
}

func init() {
	NetworkCmd.AddCommand(removeCmd)
	removeCmd.Flags().StringVarP(&label_remove, "label", "l", "", "Label for the network to be deleted")
	removeCmd.MarkFlagRequired("label")
}
