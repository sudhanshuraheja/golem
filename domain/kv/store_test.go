package kv

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
)

func TestStore(t *testing.T) {
	log := logger.NewCLILogger(6, 8)
	store := NewStore(log)

	err := store.Set("test.RandomPath", "testRandomValue")
	utils.Test().Nil(t, err)

	val, err := store.Get("test.RandomPath")
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "testRandomValue", val)

	vals, err := store.GetAll()
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "testRandomValue", vals["test.RandomPath"])

	store.Display(log, "test.RandomPath")
	err = store.Delete("test.RandomPath")
	utils.Test().Nil(t, err)

	store.SetValue("test.RandomPath", "testRandomValue2")
	store.GetValue("test.RandomPath")
	store.DeleteValue("test.RandomPath")

	err = store.Close()
	utils.Test().Nil(t, err)
}
