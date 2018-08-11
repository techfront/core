package helper

import (
	"fmt"
	"github.com/kennygrant/sanitize"
	"html/template"
)

func SanitizeHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["sanitize"] = func(s string) template.HTML {
		s, err := sanitize.HTMLAllowing(s)
		if err != nil {
			fmt.Printf("#error sanitizing html:%s", err)
			return template.HTML("")
		}
		return template.HTML(s)
	}

	return f
}
