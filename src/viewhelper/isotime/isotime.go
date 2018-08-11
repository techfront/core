package isotime

import (
	"html/template"
	"time"
)

func Helper() template.FuncMap {
	f := make(template.FuncMap)

	f["isotime"] = func(d time.Time) string {
		return d.Format("2006-01-02T15:04:05-0700")
	}

	return f
}
