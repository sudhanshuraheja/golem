package template

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/domain/kv"
	"github.com/sudhanshuraheja/golem/domain/servers"
	"github.com/sudhanshuraheja/golem/domain/vars"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

type Template struct {
	Servers servers.Servers
	Vars    vars.Vars
}

func NewTemplate(s []servers.Server, vr vars.Vars, store *kv.Store) *Template {
	t := Template{}

	t.Vars = vr
	t.trim()
	t.Servers.Merge(s)

	if store != nil {
		store, err := store.GetAll()
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
			s, err := t.Servers.Search(*servers.NewMatch(attribute, operator, value))
			if err != nil {
				return servers.Server{}
			}
			return s[0]
		},
		"match": func(attribute, operator, value string) []servers.Server {
			s, err := t.Servers.Search(*servers.NewMatch(attribute, operator, value))
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
		text = strings.ReplaceAll(text, key, v)
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
