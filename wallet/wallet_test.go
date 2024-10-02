package wallet

import (
	"math/big"
	"testing"

	"github.com/EliasManj/go-wallet/utils"
	"github.com/stretchr/testify/require"
)

func TestGetBalanceNewAccount(t *testing.T) {
	account, err := CreateNewAccount(db, "testgetbalance")
	require.NoError(t, err)
	network := Network{Label: utils.CreateAccountLabel(), ChainId: 31337, Symbol: "ETH", RpcUrl: "http://localhost:8545"}
	balance, err := GetBalance(account.Publicy, network)
	require.NoError(t, err)
	require.Equal(t, "0", balance.String())
}

func TestGetBalanceAccountWithFunds(t *testing.T) {
	// dummy account
	account := Account{
		Publicy: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Label:   utils.CreateNetworkLabel(),
	}
	err := ImportAccount(db, account)
	require.NoError(t, err)
	network := Network{Label: "test", ChainId: 31337, Symbol: "ETH", RpcUrl: "http://localhost:8545"}
	balance, err := GetBalance(account.Publicy, network)
	require.NoError(t, err)
	require.Equal(t, utils.EthToWei(big.NewFloat(10000)), balance)
}

func TestGetBalanceAfterSendingTransaction(t *testing.T) {
	// dummy accounts
	from := Account{
		Publicy:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Privatey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		Label:    utils.CreateNetworkLabel("from"),
	}
	to := Account{
		Publicy: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
		Label:   utils.CreateNetworkLabel("to"),
	}
	err := ImportAccount(db, from)
	require.NoError(t, err)
	network := Network{Label: "test", ChainId: 31337, Symbol: "ETH", RpcUrl: "http://localhost:8545"}

	amountToSend := utils.EthToWei(big.NewFloat(1))
	initialBalance, err := GetBalance(from.Publicy, network)
	require.NoError(t, err)

	tx, err := SendETH(from.Privatey, to.Publicy, amountToSend, network)
	require.NoError(t, err)
	require.Equal(t, from.Publicy, tx.From)
	require.Equal(t, to.Publicy, tx.To)
	require.Equal(t, amountToSend, tx.Amount)

	// Calculate the expected final balance
	gasUsed := tx.GasUsed   // Replace with the actual method to get gas used
	gasPrice := tx.GasPrice // Replace with the actual method to get gas price
	gasCost := new(big.Int).Mul(gasUsed, gasPrice)

	expectedFinalBalance := new(big.Int).Sub(initialBalance, tx.Amount)
	expectedFinalBalance.Sub(expectedFinalBalance, gasCost)

	finalBalance, err := GetBalance(from.Publicy, network)
	require.NoError(t, err)
	require.Equal(t, expectedFinalBalance, finalBalance)
}
