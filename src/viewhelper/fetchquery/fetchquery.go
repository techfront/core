package fetchquery

import (
	"html/template"
	"net/url"
)

func Helper() template.FuncMap {
	f := make(template.FuncMap)

	f["fetchquery"] = func(URL *url.URL) string {
		return "?" + URL.RawQuery
	}

	return f
}
