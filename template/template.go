package template

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/kv"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
	"github.com/sudhanshuraheja/golem/servers"
)

type Template struct {
	Servers []servers.Server
	Vars    map[string]string
}

func NewTemplate(s []servers.Server, c map[string]string, k *kv.KV) *Template {
	t := Template{}

	t.Vars = make(map[string]string)
	if c != nil {
		t.Vars = c
	}
	t.trim()

	t.Servers = append(t.Servers, s...)

	if k != nil {
		store, err := k.GetAll()
		if err == nil {
			for key, value := range store {
				storeKey := fmt.Sprintf("kv.%s", key)
				t.Vars[storeKey] = value
			}
		}
	}

	return &t
}

func (t *Template) Execute(text string) (string, error) {
	var err error

	text, err = t.replaceVars(text)
	if err != nil {
		return text, err
	}

	tpl := template.New("template").Funcs(template.FuncMap{
		"matchOne": func(attribute, operator, value string) servers.Server {
			s, err := servers.NewMatch(attribute, operator, value).Find(t.Servers)
			if err != nil {
				return servers.Server{}
			}
			return s[0]
		},
		"match": func(attribute, operator, value string) []servers.Server {
			s, err := servers.NewMatch(attribute, operator, value).Find(t.Servers)
			if err != nil {
				return []servers.Server{}
			}
			return s
		},
	})

	tpl, err = tpl.Parse(text)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = tpl.Execute(&b, t)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (t *Template) trim() {
	for k, v := range t.Vars {
		t.Vars[k] = strings.TrimSuffix(v, "\n")
	}
}

func (t *Template) replaceVars(text string) (string, error) {
	for k, v := range t.Vars {
		key := fmt.Sprintf("@golem.%s", k)
		text = strings.Replace(text, key, v, -1)
	}
	return text, t.checkVars(text)
}

func (t *Template) checkVars(text string) error {
	regExpKV := regexp.MustCompile(`@golem\.[\w]+\.[\w]+\.[\w]+`)
	regExpVars := regexp.MustCompile(`@golem\.[\w]+`)

	kvMatches := regExpKV.FindAllString(text, -1)
	varMatches := regExpVars.FindAllString(text, -1)
	kvMatches = append(kvMatches, varMatches...)
	kvMatches = localutils.ArrayUnique(kvMatches)

	if len(kvMatches) > 0 || len(varMatches) > 0 {
		return fmt.Errorf(
			"%s %s",
			logger.Red("found templates with no matches"),
			logger.RedBold(strings.Join(kvMatches, ", ")),
		)
	}
	return nil
}
