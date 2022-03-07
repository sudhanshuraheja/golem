package template

import (
	"testing"

	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/servers"
)

func TestTemplates(t *testing.T) {
	publicIP := "127.0.0.1"
	servers := []servers.Server{
		{
			Name:     "one",
			HostName: []string{"one"},
			Port:     22,
			Tags:     []string{"one"},
		},
		{
			Name:     "two",
			PublicIP: publicIP,
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
		Servers: servers,
	}

	message := "foo:@golem.foo {{ (matchOne \"tags\" \"contains\" \"three\").Name }}    {{ range $_, $s := (match \"tags\" \"contains\" \"one\") -}}{{- ($s).Name -}},{{- end -}}"

	_, err := tpl.Execute(message)
	utils.Test().Contains(t, err.Error(), "found templates with no matches @golem.foo")

	tpl = &Template{
		Vars:    map[string]string{"foo": "bar"},
		Servers: servers,
	}

	message = "foo:@golem.foo {{ (matchOne \"tags\" \"contains\" \"three\").Name }}    {{ range $_, $s := (match \"tags\" \"contains\" \"one\") -}}{{- ($s).Name -}},{{- end -}}"

	txt, err := tpl.Execute(message)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "foo:bar three    one,two,three,", txt)

}
