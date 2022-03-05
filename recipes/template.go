package recipes

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/sudhanshuraheja/golem/config"
)

type Template struct {
	Servers []config.Server
	Vars    map[string]string
}

func (t *Template) Trim() {
	for k, v := range t.Vars {
		t.Vars[k] = strings.TrimSuffix(v, "\n")
	}
}

func (t *Template) Execute(text string) (string, error) {
	text = t.ReplaceVars(text)
	tpl := template.New("template").Funcs(template.FuncMap{
		"matchOne": func(attribute, operator, value string) config.Server {
			s, err := NewMatch(config.Match{
				Attribute: attribute,
				Operator:  operator,
				Value:     value,
			}).Find(t.Servers)
			if err != nil {
				return config.Server{}
			}
			return s[0]
		},
		"match": func(attribute, operator, value string) []config.Server {
			s, err := NewMatch(config.Match{
				Attribute: attribute,
				Operator:  operator,
				Value:     value,
			}).Find(t.Servers)
			if err != nil {
				return []config.Server{}
			}
			return s
		},
	})

	var err error
	tpl, err = tpl.Parse(text)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = tpl.Execute(&b, t)
	if err != nil {
		return "", err
	}
	fmt.Println("==>", b.String())
	return b.String(), nil
}

func (t *Template) ReplaceVars(text string) string {
	for k, v := range t.Vars {
		key := fmt.Sprintf("@golem.%s", k)
		text = strings.Replace(text, key, v, -1)
	}
	return text
}
