package db

import (
	"fmt"
	"strconv"

	"go.etcd.io/bbolt"
)

type Db struct {
	db *bbolt.DB
}

func NewDb(path string) (*Db, error) {
	options := &bbolt.Options{}
	db, err := bbolt.Open(path, 0600, options)
	if err != nil {
		return nil, err
	}
	return &Db{db: db}, nil
}

func (d *Db) Close() error {
	return d.db.Close()
}

func (d *Db) Init() error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("nodes")); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("projects")); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("project_nodes")); err != nil {
			return err
		}
		return nil
	})
}

func (d *Db) GetNextID(bucketName string) (int, error) {
	var id int
	err := d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}

		key, _ := b.Cursor().Last()
		if key == nil {
			id = 1
		} else {
			lastID, err := strconv.Atoi(string(key))
			if err != nil {
				return err
			}
			id = lastID + 1
		}
		return nil
	})
	return id, err
}
