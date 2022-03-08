package artifacts

type Artifacts []Artifact

func (a *Artifacts) Append(art Artifact) {
	*a = append(*a, art)
}

func (a *Artifacts) Merge(arts Artifacts) {
	*a = append(*a, arts...)
}
