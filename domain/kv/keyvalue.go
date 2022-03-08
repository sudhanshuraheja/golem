package kv

type KeyValue struct {
	Path  string `hcl:"path"`
	Value string `hcl:"value"`
}
