package network

import (
	"fmt"
	"log"
	"strconv"

	"github.com/EliasManj/go-wallet/utils"
	"github.com/EliasManj/go-wallet/wallet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	label            string
	rpcURL           string
	chainID          string
	symbol           string
	blockExplorerURL string
)

// Save the network to the database
func saveNetwork(label, rpcURL, chainID, symbol, blockExplorerURL string) error {
	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		utils.WriteToBucket(db, "networks", label, []byte(fmt.Sprintf(`{"rpcURL": "%s", "chainID": "%s", "symbol": "%s", "blockExplorerURL": "%s"}`, rpcURL, chainID, symbol, blockExplorerURL)))
	}
	return nil
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new network to the wallet",
	Long:  `Add a new network to the wallet long description`,
	Run: func(cmd *cobra.Command, args []string) {

		var err error

		db, err := utils.OpenDB(viper.GetString("database_file_path"))
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()

		chainID, err := strconv.Atoi(chainID)

		network := wallet.Network{
			Label:   label,
			RpcUrl:  rpcURL,
			ChainId: chainID,
			Symbol:  symbol,
		}

		err = wallet.AddNetwork(db, network)

		if err != nil {
			log.Fatalf("Failed to save network: %v", err)
		}
		fmt.Println("Network added successfully")
	},
}

func init() {
	NetworkCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&label, "label", "l", "", "Label for the network to be identified with")
	addCmd.Flags().StringVarP(&rpcURL, "rpc-url", "r", "", "RPC URL for the network to be added")
	addCmd.Flags().StringVarP(&chainID, "chain-id", "c", "", "Chain ID for the network to be added")
	addCmd.Flags().StringVarP(&symbol, "symbol", "s", "", "Symbol for the network to be added")
	addCmd.Flags().StringVarP(&blockExplorerURL, "block-explorer-url", "b", "", "Block Explorer URL for the network to be added")
	addCmd.MarkFlagRequired("label")
	addCmd.MarkFlagRequired("rpc-url")
	addCmd.MarkFlagRequired("chain-id")
	addCmd.MarkFlagRequired("symbol")
}
