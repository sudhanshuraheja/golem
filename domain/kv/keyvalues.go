package kv

type KeyValues []*KeyValue

func (kv *KeyValues) Append(keyValue KeyValue) {
	*kv = append(*kv, &keyValue)
}

func (kv *KeyValues) Merge(keyValues KeyValues) {
	*kv = append(*kv, keyValues...)
}

func (kv *KeyValues) Setup(store *Store) (bool, error) {
	setup := false
	for _, keyValue := range *kv {
		st, err := keyValue.Setup(store)
		if err != nil {
			return setup, err
		}
		setup = st
	}
	return setup, nil
}
