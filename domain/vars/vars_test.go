package vars

import (
	"testing"

	"github.com/betas-in/utils"
)

func TestVars(t *testing.T) {
	v := NewVars()
	utils.Test().Equals(t, 0, len(*v))

	v.Add("foo", "bar")
	utils.Test().Equals(t, 1, len(*v))
}
