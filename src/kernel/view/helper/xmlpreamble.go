package helper

import (
	"html/template"
)

func XMLPreambleHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["xmlpreamble"] = func() string {
		return `<?xml version="1.0" encoding="UTF-8"?>`
	}

	return f
}
