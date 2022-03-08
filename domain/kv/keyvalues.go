package kv

type KeyValues []KeyValue

func (kv *KeyValues) Append(keyValue KeyValue) {
	*kv = append(*kv, keyValue)
}

func (kv *KeyValues) Merge(keyValues KeyValues) {
	*kv = append(*kv, keyValues...)
}
