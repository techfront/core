package helper

import (
	"html/template"
)

func EmptyHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["empty"] = func() map[string]interface{} {
		return map[string]interface{}{}
	}

	return f
}
