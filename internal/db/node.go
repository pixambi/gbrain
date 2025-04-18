package db

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go.etcd.io/bbolt"
)

type Node struct {
	ID        int
	Title     string
	Content   string
	ProjectID int
}

func (d *Db) GetNodes() ([]Node, error) {
	var nodes []Node

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("nodes"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		return b.ForEach(func(k, v []byte) error {
			var node Node
			if err := json.Unmarshal(v, &node); err != nil {
				return err
			}
			nodes = append(nodes, node)
			return nil
		})
	})
	return nodes, err
}

func (d *Db) GetNodesByProjectID(projectID int) ([]Node, error) {
	var nodes []Node

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("nodes"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var node Node
			if err := json.Unmarshal(v, &node); err != nil {
				return err
			}
			if node.ProjectID == projectID {
				nodes = append(nodes, node)
			}
			return nil
		})
	})
	return nodes, err
}

func (d *Db) GetNode(id int) (Node, error) {
	var node Node

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("nodes"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		v := b.Get([]byte(strconv.Itoa(id)))
		if v == nil {
			return fmt.Errorf("node not found")
		}

		return json.Unmarshal(v, &node)
	})
	return node, err
}

func (d *Db) AddNode(node Node) error {
	if node.ID == 0 {
		id, err := d.GetNextID("nodes")
		if err != nil {
			return err
		}
		node.ID = id
	}

	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("nodes"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		buf, err := json.Marshal(node)
		if err != nil {
			return err
		}

		key := []byte(strconv.Itoa(node.ID))
		return b.Put(key, buf)
	})
}

func (d *Db) UpdateNode(node Node) error {
	return d.AddNode(node)
}

func (d *Db) DeleteNode(id int) error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("nodes"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		return b.Delete([]byte(strconv.Itoa(id)))
	})
}
