package helper

import (
	"html/template"
)

func EscapeHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["escape"] = func(s string) string {
		return template.HTMLEscapeString(s)
	}

	return f
}
