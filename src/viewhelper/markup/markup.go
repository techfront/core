package markup

import (
	"github.com/kennygrant/sanitize"
	"html/template"
	"strings"
)

func Helper() template.FuncMap {
	f := make(template.FuncMap)

	f["markup"] = func(s string) template.HTML {
		s = strings.Replace(s, "\n", "</p><p>", -1)

		s, err := sanitize.HTMLAllowing(s)
		if err != nil {
			return template.HTML("")
		}

		return template.HTML(s)
	}

	return f
}
