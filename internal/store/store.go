package store

import (
	"encoding/json"
	"errors"

	"github.com/ostafen/kronos/internal/model"
	bolt "go.etcd.io/bbolt"
)

var ErrNotPresent = errors.New("no such schedule")

type ScheduleStore interface {
	Get(id string) (*model.Schedule, error)
	InsertOrUpdate(sched *model.Schedule) error
	Delete(id string) error
	Iterate(iterFunc func(sched *model.Schedule) error) error
}

var (
	bucketName = []byte("schedules")
)

type boltStore struct {
	db *bolt.DB
}

func New(path string) (ScheduleStore, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	s := &boltStore{
		db: db,
	}
	return s, s.init()
}

func (s *boltStore) init() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

func (s *boltStore) Get(id string) (*model.Schedule, error) {
	var sched model.Schedule
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		data := b.Get([]byte(id))
		if data == nil {
			return ErrNotPresent
		}
		return json.Unmarshal(data, &sched)
	})
	return &sched, err
}

func (s *boltStore) InsertOrUpdate(sched *model.Schedule) error {
	data, err := json.Marshal(sched)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Put([]byte(sched.ID), data)
	})
}

func (s *boltStore) Delete(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Delete([]byte(id))
	})
}

func (s *boltStore) Iterate(iterFunc func(sched *model.Schedule) error) error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		cursor := b.Cursor()

		for k, v := cursor.First(); k != nil && v != nil; k, v = cursor.Next() {
			var sched model.Schedule
			if err := json.Unmarshal(v, &sched); err != nil {
				return err
			}

			if err := iterFunc(&sched); err != nil {
				return err
			}
		}
		return nil
	})
}
