package db

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

func (d *Db) GetNodeByTitle(title string, projectID int) (Node, error) {
	var node Node
	var found bool

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("nodes"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var n Node
			if err := json.Unmarshal(v, &n); err != nil {
				return err
			}
			if n.Title == title && n.ProjectID == projectID {
				node = n
				found = true
			}
			return nil
		})
	})

	if err != nil {
		return Node{}, err
	}

	if !found {
		return Node{}, fmt.Errorf("node not found")
	}

	return node, nil
}
