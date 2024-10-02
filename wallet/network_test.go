// main_test.go
package wallet

import (
	"log"
	"os"
	"testing"

	"github.com/EliasManj/go-wallet/utils"
	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

func TestMain(m *testing.M) {
	// Setup: Open a temporary Bolt database
	var err error
	db, err = bolt.Open("test.db", 0600, nil)
	if err != nil {
		log.Fatalf("Failed to open test database: %s", err)
	}

	code := m.Run()

	db.Close()
	os.Remove("test.db")

	os.Exit(code)
}

func TestAddNetworkWithSameLabel(t *testing.T) {
	// add network
	network := Network{Label: "same", ChainId: 123, Symbol: "ETH", RpcUrl: "http://localhost:8545"}
	if err := AddNetwork(db, network); err != nil {
		t.Fatalf("Failed to add network: %s", err)
	}
	// add network with the same label
	if err := AddNetwork(db, network); err == nil {
		t.Fatalf("Expected error when adding network with the same label")
	}
}

func TestAddAndDeleteNetwork(t *testing.T) {
	// add network
	network := Network{Label: "delete", ChainId: 123, Symbol: "ETH", RpcUrl: "http://localhost:8545"}
	if err := AddNetwork(db, network); err != nil {
		t.Fatalf("Failed to add network: %s", err)
	}
	// delete the network
	if err := deleteNetwork(db, network.Label); err != nil {
		t.Fatalf("Failed to delete network: %s", err)
	}
}

func TestAddAndListNetworks(t *testing.T) {
	// Test data
	testNetworks := []Network{
		{Label: utils.CreateAccountLabel("net1"), ChainId: 123, Symbol: "ETH", RpcUrl: "http://localhost:8545"},
		{Label: utils.CreateAccountLabel("net2"), ChainId: 456, Symbol: "BTC", RpcUrl: "http://localhost:8545"},
	}

	// Test AddNetwork function
	for _, network := range testNetworks {
		if err := AddNetwork(db, network); err != nil {
			t.Fatalf("Failed to add network: %s", err)
		}
	}

	// Test listNetworks function
	networks, err := ListNetworks(db)
	if err != nil {
		t.Fatalf("Failed to list networks: %s", err)
	}

	// Verify that the networks added are present
	if len(networks) != len(testNetworks) {
		t.Fatalf("Expected %d networks, got %d", len(testNetworks), len(networks))
	}

	for _, network := range testNetworks {
		found := false
		for _, net := range networks {
			if net.Label == network.Label && net.ChainId == network.ChainId && net.Symbol == network.Symbol {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Network %v not found in listed networks", network)
		}
	}
}

func TestAddNetworkWithInvalidRpcUrl(t *testing.T) {
	// add network with invalid rpcUrl
	network := Network{Label: "invalid", ChainId: 123, Symbol: "ETH", RpcUrl: "localhost:8545"}
	if err := AddNetwork(db, network); err == nil {
		t.Fatalf("Expected error when adding network with invalid rpcUrl")
	}
}
