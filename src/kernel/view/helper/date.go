package helper

import (
	"fmt"
	"html/template"
	"time"
)

func DateHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["date"] = func(t time.Time, formats ...string) template.HTML {

		layout := "Jan 2, 2006"
		if len(formats) > 0 {
			layout = formats[0]
		}
		value := fmt.Sprintf(t.Format(layout))

		return template.HTML(template.HTMLEscapeString(value))
	}

	return f
}
