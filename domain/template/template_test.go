package template

import (
	"testing"

	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
	"github.com/sudhanshuraheja/golem/domain/kv"
	"github.com/sudhanshuraheja/golem/domain/servers"
)

func TestTemplates(t *testing.T) {
	publicIP := "127.0.0.1"
	servers := []servers.Server{
		{
			Name:     "one",
			HostName: &[]string{"one"},
			Port:     22,
			Tags:     &[]string{"one"},
		},
		{
			Name:     "two",
			PublicIP: &publicIP,
			HostName: &[]string{"two"},
			Port:     22,
			Tags:     &[]string{"one", "two"},
		},
		{
			Name:     "three",
			HostName: &[]string{"three"},
			Port:     22,
			Tags:     &[]string{"one", "two", "three"},
		},
	}

	message := `
foo:@golem.foo
FOO:@golem.FOO
{{- (matchOne "tags" "contains" "one").Name -}}
{{- if $s := (matchOne "tag" "contains" "one") -}}
	{{- ($s).Name -}},
{{- end -}}
{{- range $_, $s := (match "tag" "contains" "one") -}}
	{{- ($s).Name -}},
{{- end -}}
{{- range $_, $s := (match "tags" "contains" "one") -}}
	{{- ($s).Name -}},
{{- end -}}`
	vars := make(map[string]string)
	log := logger.NewCLILogger(6, 8)
	store := kv.NewStore(log)
	err := store.Set("test.random_path", "random2")
	utils.Test().Nil(t, err)

	tpl := NewTemplate(servers, nil, nil)
	_, err = tpl.Execute(message)
	utils.Test().Contains(t, err.Error(), "@golem.foo")
	utils.Test().Contains(t, err.Error(), "@golem.FOO")

	vars["foo"] = "bar"
	vars["FOO"] = "BAR"
	tpl = NewTemplate(servers, vars, nil)
	txt, err := tpl.Execute(message)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "\nfoo:bar\nFOO:BARone,one,two,three,", txt)

	message += `PATH:@golem.kv.test.random_path`
	tpl = NewTemplate(servers, vars, store)
	txt, err = tpl.Execute(message)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, "\nfoo:bar\nFOO:BARone,one,two,three,PATH:random2", txt)

	err = store.Delete("test.random_path")
	utils.Test().Nil(t, err)
}
