package helper

import (
	"html/template"
)

func EscapeURLHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["escapeurl"] = func(s string) string {
		return template.URLQueryEscaper(s)
	}

	return f
}
