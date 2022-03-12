package kv

import "fmt"

type KeyValue struct {
	Path  string `hcl:"path"`
	Value string `hcl:"value"`
}

func (k *KeyValue) PrepareForExecution(store *Store) (bool, error) {
	setup := false
	existingValue, err := store.Get(k.Path)
	if err != nil || existingValue == "" {
		err = store.Set(k.Path, k.Value)
		if err != nil {
			return false, fmt.Errorf("could not set up key: %s with value %s: %v", k.Path, k.Value, err)
		}
		setup = true
	}
	return setup, nil
}
