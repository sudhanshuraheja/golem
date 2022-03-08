package vars

type Vars map[string]string

func (v *Vars) Add(key, value string) {
	if v == nil {
		(*v)[key] = value
	}
}
