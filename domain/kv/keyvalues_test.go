package kv

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
)

func TestKeyValues(t *testing.T) {
	kv1 := KeyValue{Path: "test.key1", Value: "value1"}
	kv2 := KeyValue{Path: "test.key2", Value: "value2"}

	kvs1 := KeyValues{}
	kvs1.Append(kv1)
	utils.Test().Equals(t, 1, len(kvs1))

	kvs2 := KeyValues{}
	kvs2.Append(kv2)
	utils.Test().Equals(t, 1, len(kvs2))

	kvs1.Merge(kvs2)
	utils.Test().Equals(t, 2, len(kvs1))

	log := logger.NewCLILogger(6, 8)
	store := NewStore(log)
	done, err := kvs1.PrepareForExecution(store)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, done)

	err = store.Delete(kv1.Path)
	utils.Test().Nil(t, err)
	err = store.Delete(kv2.Path)
	utils.Test().Nil(t, err)

	err = store.Close()
	utils.Test().Nil(t, err)
}
