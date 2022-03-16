package keyvalue

import (
	"github.com/sudhanshuraheja/golem/domain/kv"
	"github.com/sudhanshuraheja/golem/domain/template"
)

type KeyValues []*KeyValue

func (kv *KeyValues) Append(keyValue KeyValue) {
	*kv = append(*kv, &keyValue)
}

func (kv *KeyValues) Merge(keyValues KeyValues) {
	*kv = append(*kv, keyValues...)
}

func (kv *KeyValues) PrepareForExecution(store *kv.Store, tpl *template.Template) (bool, error) {
	setup := false
	for _, keyValue := range *kv {
		st, err := keyValue.PrepareForExecution(store, tpl)
		if err != nil {
			return setup, err
		}
		setup = st
	}
	return setup, nil
}
