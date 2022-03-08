package kv

import (
	"fmt"
)

type KeyValues []KeyValue

func (kv *KeyValues) Append(keyValue KeyValue) {
	*kv = append(*kv, keyValue)
}

func (kv *KeyValues) Merge(keyValues KeyValues) {
	*kv = append(*kv, keyValues...)
}

func (kv *KeyValues) Setup(store *Store) (bool, error) {
	setup := false
	for _, keyValue := range *kv {
		key := keyValue.Path
		value := keyValue.Value

		existingValue, err := store.Get(key)
		if err != nil || existingValue == "" {
			err = store.Set(key, value)
			if err != nil {
				return false, fmt.Errorf("could not set up key: %s with value %s: %v", key, value, err)
			}
			setup = true
		}
	}
	return setup, nil
}
