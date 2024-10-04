package wallet

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
)

type Network struct {
	Label    string `json:"label"`
	ChainId  int    `json:"chainId"`
	Symbol   string `json:"symbol"`
	RpcUrl   string `json:"rpcUrl"`
	Selected bool   `json:"selected"`
}

// addNetwork adds a new network to the Bolt database
func AddNetwork(db *bolt.DB, network Network) error {
	// Validate the RpcUrl
	if !strings.HasPrefix(network.RpcUrl, "http://") && !strings.HasPrefix(network.RpcUrl, "https://") {
		return fmt.Errorf("rpcUrl must start with 'http://' or 'https://'")
	}
	network.Selected = true
	// Open a writable transaction
	err := db.Update(func(tx *bolt.Tx) error {
		// Create or get the bucket "networks"
		bucket, err := tx.CreateBucketIfNotExists([]byte("networks"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		// Check if the network with the same label already exists
		existing := bucket.Get([]byte(network.Label))
		if existing != nil {
			return fmt.Errorf("network with label %s already exists", network.Label)
		}

		// Marshal the network object to JSON
		networkJSON, err := json.Marshal(network)
		if err != nil {
			return fmt.Errorf("json marshal: %s", err)
		}

		// Use the label as the key and store the JSON as the value
		return bucket.Put([]byte(network.Label), networkJSON)
	})

	if err != nil {
		return err
	}

	return SelectNetwork(db, network.Label)
}

// listNetworks retrieves all networks from the Bolt database
func ListNetworks(db *bolt.DB) ([]Network, error) {
	var networks []Network

	// Open a read-only transaction
	err := db.View(func(tx *bolt.Tx) error {
		// Get the bucket "networks"
		bucket := tx.Bucket([]byte("networks"))
		if bucket == nil {
			return nil
		}

		// Iterate through all key-value pairs in the bucket
		return bucket.ForEach(func(k, v []byte) error {
			var network Network
			// Unmarshal the JSON value to a Network object
			err := json.Unmarshal(v, &network)
			if err != nil {
				return fmt.Errorf("json unmarshal: %s", err)
			}
			// Append to the list of networks
			networks = append(networks, network)
			return nil
		})
	})

	return networks, err
}

// deleteNetwork removes a network from the Bolt database
func deleteNetwork(db *bolt.DB, label string) error {
	// Open a writable transaction
	err := db.Update(func(tx *bolt.Tx) error {
		// Get the bucket "networks"
		bucket := tx.Bucket([]byte("networks"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		// Delete the key-value pair with the given label
		return bucket.Delete([]byte(label))
	})

	return err
}

func UpdateNetwork(db *bolt.DB, network Network) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("networks"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		networkJSON, err := json.Marshal(network)
		if err != nil {
			return fmt.Errorf("json marshal: %s", err)
		}

		return bucket.Put([]byte(network.Label), networkJSON)
	})

	return err
}

func GetNetwork(db *bolt.DB, label string) (Network, error) {
	var network Network

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("networks"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		networkJSON := bucket.Get([]byte(label))
		if networkJSON == nil {
			return fmt.Errorf("network not found")
		}

		err := json.Unmarshal(networkJSON, &network)
		if err != nil {
			return fmt.Errorf("json unmarshal: %s", err)
		}

		return nil
	})

	return network, err
}

func SelectNetwork(db *bolt.DB, label string) error {
	// deselect currente selected network
	selectedNetwork, err := GetSelectedNetwork(db)
	if err == nil {
		selectedNetwork.Selected = false
		err = UpdateNetwork(db, selectedNetwork)
		if err != nil {
			return err
		}
	}

	network, err := GetNetwork(db, label)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("selected"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return bucket.Put([]byte("network"), []byte(label))
	})
	if err != nil {
		return err
	}

	network.Selected = true
	return UpdateNetwork(db, network)
}

func GetSelectedNetwork(db *bolt.DB) (Network, error) {
	var network Network

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("selected"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		label := bucket.Get([]byte("network"))
		if label == nil {
			return fmt.Errorf("network not selected")
		}

		networkJSON := tx.Bucket([]byte("networks")).Get(label)
		if networkJSON == nil {
			return fmt.Errorf("network not found")
		}

		err := json.Unmarshal(networkJSON, &network)
		if err != nil {
			return fmt.Errorf("json unmarshal: %s", err)
		}

		return nil
	})

	return network, err
}
