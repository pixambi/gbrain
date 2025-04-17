package db

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go.etcd.io/bbolt"
)

type Project struct {
	ID   int
	Name string
}

func (d *Db) GetProjects() ([]Project, error) {
	var projects []Project

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("projects"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var project Project
			if err := json.Unmarshal(v, &project); err != nil {
				return err
			}
			projects = append(projects, project)
			return nil
		})
	})
	return projects, err
}

func (d *Db) GetProject(id int) (Project, error) {
	var project Project

	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("projects"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		v := b.Get([]byte(strconv.Itoa(id)))
		if v == nil {
			return fmt.Errorf("project not found")
		}

		return json.Unmarshal(v, &project)
	})
	return project, err
}

func (d *Db) AddProject(project Project) error {
	if project.ID == 0 {
		id, err := d.GetNextID("projects")
		if err != nil {
			return err
		}
		project.ID = id
	}

	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("projects"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		buf, err := json.Marshal(project)
		if err != nil {
			return err
		}

		key := []byte(strconv.Itoa(project.ID))
		return b.Put(key, buf)
	})
}

func (d *Db) UpdateProject(project Project) error {
	return d.AddProject(project)
}

func (d *Db) DeleteProject(id int) error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("projects"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		nodes, err := d.GetNodesByProjectID(id)
		if err != nil {
			return err
		}

		for _, node := range nodes {
			if err := d.DeleteNode(node.ID); err != nil {
				return err
			}
		}

		return b.Delete([]byte(strconv.Itoa(id)))
	})
}
