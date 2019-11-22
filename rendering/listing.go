package rendering

import (
	"bytes"
	"html/template"
)

func Render(template *template.Template, data interface{}) (rendered string, err error) {
	writer := new(bytes.Buffer)
	err = template.Execute(writer, data)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}
