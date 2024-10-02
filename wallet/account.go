package wallet

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	Label    string `json:"label"`
	Publicy  string `json:"pubic"`
	Privatey string `json:"private"`
	Tokens   []string
	Selected bool
}

func GenerateKeyPair() (string, string, error) {
	// Generate a new private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	// Get the public key
	publicKey := privateKey.PublicKey

	// Convert public key to Ethereum address
	address := crypto.PubkeyToAddress(publicKey).Hex()

	privateKeyBytes := crypto.FromECDSA(privateKey)
	// Encode the byte slice to a hex string
	privateKeyHex := fmt.Sprintf("%x", privateKeyBytes)

	return privateKeyHex, address, nil
}

func CreateNewAccount(db *bolt.DB, label string) (Account, error) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		return Account{}, err
	}

	acc := Account{
		Label:    label,
		Publicy:  publicKey,
		Privatey: privateKey,
	}
	return acc, ImportAccount(db, acc)
}

func ImportAccount(db *bolt.DB, account Account) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("accounts"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		existing := bucket.Get([]byte(account.Label))
		if existing != nil {
			return fmt.Errorf("account with label %s already exists", account.Label)
		}

		accountJSON, err := json.Marshal(account)
		if err != nil {
			return fmt.Errorf("json marshal: %s", err)
		}

		return bucket.Put([]byte(account.Label), accountJSON)
	})

	return err
}

func AddTokenToAccount(db *bolt.DB, accountLabel, tokenAddress string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("accounts"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		accountJSON := bucket.Get([]byte(accountLabel))
		if accountJSON == nil {
			return fmt.Errorf("account not found")
		}

		var account Account
		err := json.Unmarshal(accountJSON, &account)
		if err != nil {
			return fmt.Errorf("json unmarshal: %s", err)
		}

		account.Tokens = append(account.Tokens, tokenAddress)

		accountJSON, err = json.Marshal(account)
		if err != nil {
			return fmt.Errorf("json marshal: %s", err)
		}

		return bucket.Put([]byte(account.Label), accountJSON)
	})

	return err
}

func GetAccount(db *bolt.DB, label string) (Account, error) {
	var account Account

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("accounts"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		accountJSON := bucket.Get([]byte(label))
		if accountJSON == nil {
			return fmt.Errorf("account not found")
		}

		err := json.Unmarshal(accountJSON, &account)
		if err != nil {
			return fmt.Errorf("json unmarshal: %s", err)
		}

		return nil
	})

	return account, err
}

func UpdateAccount(db *bolt.DB, account Account) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("accounts"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		accountJSON, err := json.Marshal(account)
		if err != nil {
			return fmt.Errorf("json marshal: %s", err)
		}

		return bucket.Put([]byte(account.Label), accountJSON)
	})

	return err
}

func RemoveAccount(db *bolt.DB, label string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("accounts"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		return bucket.Delete([]byte(label))
	})

	return err
}

func ListAccounts(db *bolt.DB) ([]Account, error) {
	var accounts []Account
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("accounts"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		err = bucket.ForEach(func(k, v []byte) error {
			var account Account
			err := json.Unmarshal(v, &account)
			if err != nil {
				return fmt.Errorf("json unmarshal: %s", err)
			}

			accounts = append(accounts, account)
			return nil
		})

		return err
	})

	return accounts, err
}

func SelectAccount(db *bolt.DB, label string) error {
	account, err := GetAccount(db, label)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("selected"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return bucket.Put([]byte("acc"), []byte(label))
	})
	if err != nil {
		return err
	}

	account.Selected = true
	return UpdateAccount(db, account)
}

func GetSelectedAccount(db *bolt.DB) (Account, error) {
	var account Account
	var label string

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("selected"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		label = string(bucket.Get([]byte("acc")))
		if label == "" {
			return fmt.Errorf("no account selected")
		}
		return nil
	})

	if err != nil {
		return account, err
	}
	account, err = GetAccount(db, label)

	if err != nil {
		return account, err
	}

	account.Selected = true

	return account, nil
}
