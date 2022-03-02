package recipes

import (
	"testing"

	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/config"
)

func TestTemplates(t *testing.T) {
	publicIP := "127.0.0.1"
	servers := []config.Server{
		{
			Name:     "one",
			HostName: []string{"one"},
			Port:     22,
			Tags:     []string{"one"},
		},
		{
			Name:     "two",
			PublicIP: &publicIP,
			HostName: []string{"two"},
			Port:     22,
			Tags:     []string{"one", "two"},
		},
		{
			Name:     "three",
			HostName: []string{"three"},
			Port:     22,
			Tags:     []string{"one", "two", "three"},
		},
	}

	tpl := &Template{
		Vars:    map[string]string{"foo": "bar"},
		Servers: servers,
	}

	message := "foo:{{ .Vars.foo}} {{ (matchOne \"tags\" \"contains\" \"three\").Name }}    {{ range $_, $s := (match \"tags\" \"contains\" \"one\") -}}{{- ($s).Name -}},{{- end -}}"

	txt, err := ParseTemplate(message, tpl)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "foo:bar three    one,two,three,", txt)
}
