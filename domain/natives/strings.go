package natives

type Strings []String

func (s *Strings) Append(str string) {
	*s = append(*s, String(str))
}

func (s *Strings) Merge(strs Strings) {
	*s = append(*s, strs...)
}
