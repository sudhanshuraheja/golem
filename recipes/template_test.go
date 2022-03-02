package recipes

import (
	"testing"

	"github.com/betas-in/utils"
)

func TestTemplates(t *testing.T) {
	tpl := &Template{
		Vars: map[string]string{"foo": "bar"},
	}

	message := "foo:{{ .Vars.foo}}"

	txt, err := ParseTemplate(message, tpl)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "foo:bar", txt)
}
