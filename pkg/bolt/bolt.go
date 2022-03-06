package bolt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"go.etcd.io/bbolt"
)

type Bolt struct {
	db *bbolt.DB
}

func NewBolt(path string) (*Bolt, error) {
	b := &Bolt{}
	err := b.Open(path)
	if err != nil {
		return b, err
	}
	return b, nil
}

func (b *Bolt) Open(path string) error {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return err
	}
	b.db = db
	return nil
}

func (b *Bolt) Close() error {
	return b.db.Close()
}

func (b *Bolt) GetDB() *bbolt.DB {
	return b.db
}

func (b *Bolt) CreateBucket(bucket []byte) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return nil
	})
}

func (b *Bolt) DeleteBucket(bucket []byte) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket(bucket)
		if bu == nil {
			// The bucket does not exist
			return nil
		}
		return tx.DeleteBucket(bucket)
	})
}

func (b *Bolt) ListBuckets() ([]string, error) {
	buckets := []string{}
	err := b.db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bbolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
	})
	return buckets, err
}

func (b *Bolt) Put(bucket, key, value []byte) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket(bucket)
		if bu == nil {
			return fmt.Errorf("bucket %s does not exist", string(bucket))
		}
		return bu.Put(key, value)
	})
}

func (b *Bolt) Get(bucket, key []byte) ([]byte, error) {
	var value []byte
	err := b.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", string(bucket))
		}
		val := b.Get(key)
		value = make([]byte, len(val))
		copy(value, val)
		return nil
	})
	return value, err
}

func (b *Bolt) Delete(bucket, key []byte) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		return b.Delete(key)
	})
}

func (b *Bolt) IntToByte(v int) []byte {
	by := make([]byte, 8)
	binary.BigEndian.PutUint64(by, uint64(v))
	return by
}

func (b *Bolt) FindWithPrefix(bucket, prefix []byte) (map[string][]byte, error) {
	values := map[string][]byte{}
	err := b.db.View(func(tx *bbolt.Tx) error {
		bu := tx.Bucket(bucket)
		if bu == nil {
			return fmt.Errorf("bucket %s does not exist", string(bucket))
		}

		cursor := bu.Cursor()
		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			var key, val []byte
			key = make([]byte, len(k))
			val = make([]byte, len(v))
			copy(val, v)
			copy(key, k)
			values[string(key)] = val
		}
		return nil
	})
	return values, err
}

func (b *Bolt) FindAll(bucket []byte) (map[string][]byte, error) {
	values := map[string][]byte{}
	err := b.db.View(func(tx *bbolt.Tx) error {
		bu := tx.Bucket(bucket)
		if bu == nil {
			return fmt.Errorf("bucket %s does not exist", string(bucket))
		}

		err := bu.ForEach(func(k, v []byte) error {
			var key, val []byte
			key = make([]byte, len(k))
			val = make([]byte, len(v))
			copy(val, v)
			copy(key, k)
			values[string(key)] = val
			return nil
		})
		return err
	})
	return values, err
}

func (b *Bolt) Backup(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return b.db.View(func(tx *bbolt.Tx) error {
		_, err := tx.WriteTo(file)
		if err != nil {
			return err
		}
		return nil
	})
}
