package utils

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// OpenDB opens a BoltDB database and returns a pointer to the DB instance.
func OpenDB(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	return db, nil
}

// ReadFromBucket reads a value from the specified bucket and key.
func ReadFromBucket(db *bolt.DB, bucketName string, key string) ([]byte, error) {
	var value []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}
		value = bucket.Get([]byte(key))
		if value == nil {
			return fmt.Errorf("key %s not found in bucket %s", key, bucketName)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read from bucket: %v", err)
	}
	return value, nil
}

// WriteToBucket writes a key-value pair to the specified bucket.
func WriteToBucket(db *bolt.DB, bucketName string, key string, value []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("failed to create or open bucket %s: %v", bucketName, err)
		}
		if err := bucket.Put([]byte(key), value); err != nil {
			return fmt.Errorf("failed to write to bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	}
	return nil
}

// CloseDB closes the BoltDB database.
func CloseDB(db *bolt.DB) error {
	err := db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database: %v", err)
	}
	return nil
}
