package helper

import (
	"fmt"
	"html/template"
	"strings"
)

func TextareaHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["textarea"] = func(label string, name string, v interface{}, args ...string) template.HTML {
		attributes := ""
		if len(args) > 0 {
			attributes = strings.Join(args, " ")
		}

		fieldTemplate := `<div class="field"><label>%s</label><textarea name="%s" %s>%v</textarea></div>`

		output := fmt.Sprintf(fieldTemplate, template.HTMLEscapeString(label), template.HTMLEscapeString(name), attributes, v)

		return template.HTML(output)
	}

	return f
}
