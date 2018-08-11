package helper

import (
	"html/template"
)

func UrlHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["url"] = func(s string) template.URL {
		return template.URL(s)
	}

	return f
}
