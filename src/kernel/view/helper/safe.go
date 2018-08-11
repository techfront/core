package helper

import (
	"html/template"
)

func SafeHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["safe"] = func(s string) template.HTML {
		return template.HTML(s)
	}

	return f
}
