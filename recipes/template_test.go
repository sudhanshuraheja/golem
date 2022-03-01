package recipes

import (
	"testing"

	"github.com/betas-in/utils"
)

func TestTemplates(t *testing.T) {
	type Todo struct {
		Vars *map[string]string
	}

	todo := Todo{
		Vars: &map[string]string{"foo": "bar"},
	}

	message := "foo:{{ .Vars.foo}}"

	txt, err := ParseTemplate(message, todo)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "foo:bar", txt)
}
