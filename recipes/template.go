package recipes

import (
	"bytes"
	"html/template"

	"github.com/sudhanshuraheja/golem/config"
)

type Template struct {
	Servers []config.Server
	Vars    map[string]string
}

func ParseTemplate(text string, tp *Template) (string, error) {
	t := template.New("template").Funcs(template.FuncMap{
		"matchOne": func(attribute, operator, value string) config.Server {
			s, err := NewMatch(config.Match{
				Attribute: attribute,
				Operator:  operator,
				Value:     value,
			}).Find(tp.Servers)
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
			}).Find(tp.Servers)
			if err != nil {
				return []config.Server{}
			}
			return s
		},
	})

	var err error
	t, err = t.Parse(text)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = t.Execute(&b, tp)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
