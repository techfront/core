package helper

import (
	"fmt"
	"html/template"
)

func StyleHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["style"] = func(name string) template.HTML {
		return template.HTML(fmt.Sprintf("<link href=\"/assets/styles/%s.css\" media=\"all\" rel=\"stylesheet\" type=\"text/css\" />", template.URLQueryEscaper(name)))
	}

	return f
}
