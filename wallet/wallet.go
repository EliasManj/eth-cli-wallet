package wallet

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERC-20 ABI for balanceOf function
const erc20ABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

type Transaction struct {
	From     string
	To       string
	Amount   *big.Int
	Network  Network
	Hash     string
	GasUsed  *big.Int
	GasPrice *big.Int
}

func GetBalance(address string, network Network) (*big.Int, error) {
	client, err := ethclient.Dial(network.RpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Convert the address to a common.Address type
	account := common.HexToAddress(address)

	// Get the balance of the account at the latest block
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}

	return balance, nil
}

func GetTokenBalance(tokenAddress string, ownerAddress string, network Network) (*big.Int, error) {
	client, err := ethclient.Dial(network.RpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Convert the addresses to common.Address type
	token := common.HexToAddress(tokenAddress)
	owner := common.HexToAddress(ownerAddress)

	// Parse the ERC-20 token ABI
	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ERC-20 ABI: %v", err)
	}

	// Create a call message to query the balanceOf method
	callData, err := parsedABI.Pack("balanceOf", owner)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data for balanceOf call: %v", err)
	}

	// Perform the call
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &token,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}

	// Unpack the result into a big.Int
	var balance *big.Int
	err = parsedABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack balance: %v", err)
	}

	return balance, nil
}

// SendETH sends Ether from one account to another.
func SendETH(fromPrivateKey string, toAddress string, amount *big.Int, network Network) (Transaction, error) {
	client, err := ethclient.Dial(network.RpcUrl)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Remove "0x" prefix from the private key
	if len(fromPrivateKey) > 2 && fromPrivateKey[:2] == "0x" {
		fromPrivateKey = fromPrivateKey[2:]
	}

	// Convert the private key string to ecdsa.PrivateKey
	privateKey, err := crypto.HexToECDSA(fromPrivateKey)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid private key: %v", err)
	}

	// Get the public address of the sender
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to get nonce: %v", err)
	}

	// Set gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to get gas price: %v", err)
	}

	// Set gas limit
	gasLimit := uint64(21000) // basic transaction

	// Create the transaction
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), amount, gasLimit, gasPrice, nil)

	// Sign the transaction with the sender's private key
	chainID := big.NewInt(int64(network.ChainId))
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to send transaction: %v", err)
	}

	// Retry mechanism for fetching the receipt
	var receipt *types.Receipt
	for {
		receipt, err = client.TransactionReceipt(context.Background(), signedTx.Hash())
		if err == nil {
			break
		}
		if err.Error() == "not found" {
			time.Sleep(time.Second * 1)
			continue
		}
		return Transaction{}, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	// Return the transaction details including gas used and gas price
	return Transaction{
		From:     fromAddress.String(),
		To:       toAddress,
		Amount:   amount,
		Network:  network,
		Hash:     signedTx.Hash().Hex(),
		GasUsed:  new(big.Int).SetUint64(receipt.GasUsed),
		GasPrice: gasPrice,
	}, nil
}

func GetAddressFromPrivateKey(privateKey string) (string, error) {
	// Remove "0x" prefix from the private key
	if len(privateKey) > 2 && privateKey[:2] == "0x" {
		privateKey = privateKey[2:]
	}

	// Convert the private key string to ecdsa.PrivateKey
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// Get the public address of the sender
	fromAddress := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	return fromAddress.String(), nil
}
