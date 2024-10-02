package utils

import (
	"math/big"
	"runtime"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

func GenerateRandomString(length int) string {
	// Define the character set
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Seed the random number generator
	rand.Seed(uint64(time.Now().UnixNano()))

	// Generate the random string
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
func CreateAccountLabel(extras ...string) string {
	// Retrieve the caller's program counter, file, and line number
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	// Retrieve the function's details using the program counter
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	// Get the full function name, which includes the package path
	fullName := fn.Name()

	// Extract the function name by splitting on '/' and '.'
	parts := strings.Split(fullName, "/")
	lastPart := parts[len(parts)-1]
	functionName := strings.Split(lastPart, ".")[1]

	// Join the extras with an underscore
	extrasStr := strings.Join(extras, "_")

	// Construct the final label
	label := "account"
	if extrasStr != "" {
		label += "_" + extrasStr
	}

	return strings.ToLower(label + functionName)
}

func CreateNetworkLabel(extras ...string) string {
	// Retrieve the caller's program counter, file, and line number
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	// Retrieve the function's details using the program counter
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	// Get the full function name, which includes the package path
	fullName := fn.Name()

	// Extract the function name by splitting on '/' and '.'
	parts := strings.Split(fullName, "/")
	lastPart := parts[len(parts)-1]
	functionName := strings.Split(lastPart, ".")[1]

	// Join the extras with an underscore
	extrasStr := strings.Join(extras, "_")

	// Construct the final label
	label := "network"
	if extrasStr != "" {
		label += "_" + extrasStr
	}

	return strings.ToLower(label + functionName)
}

// WeiToEth converts Wei to ETH
func WeiToEth(wei *big.Int) *big.Float {
	eth := new(big.Float).SetInt(wei)
	eth = new(big.Float).Quo(eth, big.NewFloat(1e18))
	return eth
}

// EthToWei converts ETH to Wei
func EthToWei(eth *big.Float) *big.Int {
	// Define the conversion factor (1 ETH = 10^18 Wei)
	weiConversionFactor := big.NewFloat(1e18)

	// Multiply the ETH value by the conversion factor
	wei := new(big.Float).Mul(eth, weiConversionFactor)

	// Convert the result to *big.Int
	weiInt := new(big.Int)
	wei.Int(weiInt) // Extract the integer part of the result
	return weiInt
}
