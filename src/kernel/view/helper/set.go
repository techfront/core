package helper

import (
	"html/template"
)

func SetHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["set"] = func(m map[string]interface{}, k string, v interface{}) string {
		m[k] = v
		return ""
	}

	return f
}
