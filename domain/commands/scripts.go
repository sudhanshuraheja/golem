package commands

type Scripts []Script

func (s *Scripts) Append(scr Script) {
	*s = append(*s, scr)
}

func (s *Scripts) Merge(scrs Scripts) {
	*s = append(*s, scrs...)
}
