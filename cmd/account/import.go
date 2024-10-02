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
	import_label string
	import_key   string
)

func importAccount() {
	db, err := utils.OpenDB(viper.GetString("database_file_path"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	_, err = wallet.GetAccount(db, import_label)
	if err == nil {
		fmt.Printf("Account with label %s already exists\n", import_label)
		return
	}

	public, err := wallet.GetAddressFromPrivateKey(import_key)
	if err != nil {
		fmt.Printf("Failed to get public key from private key: %v\n", err)
		return
	}

	account := wallet.Account{
		Label:    import_label,
		Privatey: import_key,
		Publicy:  public,
	}

	err = wallet.ImportAccount(db, account)
	if err != nil {
		fmt.Printf("Failed to import account: %v\n", err)
		return
	}

	fmt.Printf("Account with label %s imported successfully\n", import_label)
}

var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "This import all the networks available for the wallet",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		importAccount()
	},
}

func init() {
	AccountCmd.AddCommand(ImportCmd)
	ImportCmd.Flags().StringVarP(&import_label, "label", "l", "", "Label for the network to be identified with")
	ImportCmd.MarkFlagRequired("label")
	ImportCmd.Flags().StringVarP(&import_key, "key", "k", "", "Label for the network to be identified with")
	ImportCmd.MarkFlagRequired("key")
}
