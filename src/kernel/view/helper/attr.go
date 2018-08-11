package helper

import (
	"html/template"
)

func AttrHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["attr"] = func(s string) template.HTMLAttr {
		return template.HTMLAttr(s)
	}

	return f
}
