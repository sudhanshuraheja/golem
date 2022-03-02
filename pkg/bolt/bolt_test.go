package bolt

import (
	"os"
	"testing"

	"github.com/betas-in/utils"
)

func TestBolt(t *testing.T) {
	path := "bolt.db"

	b, err := NewBolt(path)
	utils.Test().Nil(t, err)

	bucket := []byte("bucket1")
	key1 := []byte("secrets/golem/1")
	value1 := []byte("this is a super secret - 1")
	key2 := []byte("secrets/golem/2")
	value2 := []byte("this is a super secret - 2")
	key3 := []byte("kv/golem/3")
	value3 := []byte("this is an env var - 3")

	err = b.CreateBucket(bucket)
	utils.Test().Nil(t, err)

	err = b.DeleteBucket(bucket)
	utils.Test().Nil(t, err)

	err = b.Put(bucket, key1, value1)
	utils.Test().Contains(t, err.Error(), "does not exist")

	err = b.CreateBucket(bucket)
	utils.Test().Nil(t, err)

	err = b.Put(bucket, key1, value1)
	utils.Test().Nil(t, err)

	err = b.Put(bucket, key2, value2)
	utils.Test().Nil(t, err)

	err = b.Put(bucket, key3, value3)
	utils.Test().Nil(t, err)

	rValue1, err := b.Get(bucket, key1)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, value1, rValue1)

	rValue2, err := b.Get(bucket, key2)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, value2, rValue2)

	rValue3, err := b.Get(bucket, key3)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, value3, rValue3)

	err = b.Delete(bucket, key2)
	utils.Test().Nil(t, err)

	err = b.Put(bucket, key2, value2)
	utils.Test().Nil(t, err)

	vals, err := b.FindWithPrefix(bucket, []byte("secrets/golem/"))
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 2, len(vals))
	utils.Test().Equals(t, value1, vals[string(key1)])
	utils.Test().Equals(t, value2, vals[string(key2)])

	vals, err = b.FindAll(bucket)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, 3, len(vals))
	utils.Test().Equals(t, value1, vals[string(key1)])
	utils.Test().Equals(t, value2, vals[string(key2)])
	utils.Test().Equals(t, value3, vals[string(key3)])

	err = b.DeleteBucket(bucket)
	utils.Test().Nil(t, err)

	b.Close()

	os.Remove(path)
}
