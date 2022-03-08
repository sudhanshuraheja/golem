package vars

type Vars map[string]string

func NewVars() *Vars {
	v := Vars(make(map[string]string))
	return &v
}

func (v *Vars) Add(key, value string) {
	if v != nil {
		(*v)[key] = value
	}
}
