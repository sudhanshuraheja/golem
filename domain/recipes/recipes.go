package recipes

type Recipes []Recipe

func (r *Recipes) Append(rcp Recipe) {
	if r != nil {
		*r = append(*r, rcp)
	}
}

func (r *Recipes) Merge(rcps Recipes) {
	if r != nil {
		*r = append(*r, rcps...)
	}
}
