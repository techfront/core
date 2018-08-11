package helper

import (
	"fmt"
	"html/template"
	"strings"
)

func FieldHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["field"] = func(label string, name string, v interface{}, args ...string) template.HTML {
		attributes := ""
		if len(args) > 0 {
			attributes = strings.Join(args, " ")
		}

		if !strings.Contains(attributes, "type=") {
			attributes = attributes + " type=\"text\""
		}

		tmpl := `<div class="field"><label>%s</label><input name="%s" value="%s" %s></div>`

		if label == "" {
			tmpl = `%s<input name="%s" value="%s" %s>`
		}

		output := fmt.Sprintf(tmpl, template.HTMLEscapeString(label), template.HTMLEscapeString(name), template.HTMLEscapeString(fmt.Sprintf("%v", v)), attributes)

		return template.HTML(output)
	}

	return f
}
