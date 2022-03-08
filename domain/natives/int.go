package natives

type Int int

type Ints []Int

func (i *Ints) Append(in int) {
	*i = append(*i, Int(in))
}

func (i *Ints) Merge(ins Ints) {
	*i = append(*i, ins...)
}
