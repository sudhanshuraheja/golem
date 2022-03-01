package recipes

import (
	"bytes"
	"html/template"
)

func ParseTemplate(text string, tp interface{}) (string, error) {
	t := template.New("template")

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
