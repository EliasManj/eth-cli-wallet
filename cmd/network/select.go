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
	label_select string
)

func selectNetwork() {
	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	network, err := wallet.GetNetwork(db, label_select)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	err = wallet.SelectNetwork(db, network.Label)
	if err != nil {
		log.Fatalf("Failed to select network: %v", err)
	}

	fmt.Printf("Network %s selected\n", network.Label)
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "This selects the specified network to be the active network",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		selectNetwork()
	},
}

func init() {
	NetworkCmd.AddCommand(selectCmd)
	selectCmd.Flags().StringVarP(&label_select, "label", "l", "", "Label for the network to be identified with")
	selectCmd.MarkFlagRequired("label")
}
