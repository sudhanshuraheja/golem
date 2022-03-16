package keyvalue

import (
	"fmt"

	"github.com/sudhanshuraheja/golem/domain/kv"
	"github.com/sudhanshuraheja/golem/domain/template"
)

type KeyValue struct {
	Path  string `hcl:"path"`
	Value string `hcl:"value"`
}

func (k *KeyValue) PrepareForExecution(store *kv.Store, tpl *template.Template) (bool, error) {
	setup := false

	if k.Path != "" {
		var err error
		k.Path, err = tpl.Execute(k.Path)
		if err != nil {
			return setup, err
		}
	}

	if k.Value != "" {
		var err error
		k.Path, err = tpl.Execute(k.Path)
		if err != nil {
			return setup, err
		}
	}

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
