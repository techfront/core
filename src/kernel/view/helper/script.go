package helper

import (
	"fmt"
	"html/template"
)

func ScriptHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["script"] = func(name string) template.HTML {
		return template.HTML(fmt.Sprintf("<script src=\"/assets/scripts/%s.js\" type=\"text/javascript\"></script>", template.URLQueryEscaper(name)))
	}

	return f
}
