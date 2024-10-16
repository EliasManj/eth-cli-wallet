// main_test.go
package wallet

import (
	"log"
	"os"
	"testing"

	"github.com/EliasManj/go-wallet/utils"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"
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

func TestAddNetwork(t *testing.T) {
	network := Network{Label: "test", ChainId: 123, Symbol: "ETH", RpcUrl: "http://localhost:8545"}
	if err := AddNetwork(db, network); err != nil {
		t.Fatalf("Failed to add network: %s", err)
	}
	net, err := GetNetwork(db, network.Label)
	require.NoError(t, err)
	require.Equal(t, network.Label, net.Label)
	require.Equal(t, network.ChainId, net.ChainId)
	require.Equal(t, network.Symbol, net.Symbol)
	require.Equal(t, network.RpcUrl, net.RpcUrl)
	require.True(t, net.Selected)
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
	if err := DeleteNetwork(db, network.Label); err != nil {
		t.Fatalf("Failed to delete network: %s", err)
	}
}

func TestAddAndListNetworks(t *testing.T) {

	todel, err := ListNetworks(db)
	require.NoError(t, err)

	for _, network := range todel {
		err := DeleteNetwork(db, network.Label)
		require.NoError(t, err)
	}

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

func TestOnlyOneSelectedNetwork(t *testing.T) {
	net1 := Network{Label: utils.CreateAccountLabel("testonly1sel1"), ChainId: 123, Symbol: "ETH", RpcUrl: "http://localhost:8545", Selected: true}
	net2 := Network{Label: utils.CreateAccountLabel("testonly1sel2"), ChainId: 456, Symbol: "BTC", RpcUrl: "http://localhost:8545", Selected: true}
	net3 := Network{Label: utils.CreateAccountLabel("testonly1sel3"), ChainId: 456, Symbol: "BTC", RpcUrl: "http://localhost:8545", Selected: true}

	require.NoError(t, AddNetwork(db, net1))
	require.NoError(t, AddNetwork(db, net2))
	require.NoError(t, AddNetwork(db, net3))

	SelectNetwork(db, net1.Label)
	SelectNetwork(db, net2.Label)
	SelectNetwork(db, net3.Label)

	n1, err := GetNetwork(db, net1.Label)
	require.NoError(t, err)
	require.Equal(t, net1.Label, n1.Label)
	n2, err := GetNetwork(db, net2.Label)
	require.NoError(t, err)
	require.Equal(t, net2.Label, n2.Label)
	n3, err := GetNetwork(db, net3.Label)
	require.NoError(t, err)
	require.Equal(t, net3.Label, n3.Label)

	require.Equal(t, false, n1.Selected)
	require.Equal(t, false, n2.Selected)
	require.Equal(t, true, n3.Selected)

	selected, err := GetSelectedNetwork(db)
	require.NoError(t, err)
	require.Equal(t, net3.Label, selected.Label)
}

func TestAddNetworkAndGetSelected(t *testing.T) {
	network := Network{Label: "addandtestselected", ChainId: 123, Symbol: "ETH", RpcUrl: "http://localhost:8545"}
	require.NoError(t, AddNetwork(db, network))

	net, err := GetSelectedNetwork(db)
	require.NoError(t, err)
	require.Equal(t, network.Label, net.Label)
	require.Equal(t, network.ChainId, net.ChainId)
	require.Equal(t, network.Symbol, net.Symbol)
	require.Equal(t, network.RpcUrl, net.RpcUrl)
	require.True(t, net.Selected)
}

func TestNoNetworksToList(t *testing.T) {
	networks, err := ListNetworks(db)
	require.NoError(t, err)

	for _, network := range networks {
		err := DeleteNetwork(db, network.Label)
		require.NoError(t, err)
	}

	_, err = ListNetworks(db)
	require.NoError(t, err)

}
