package displayfloat

import (
	"fmt"
	"html/template"
)

func Helper() template.FuncMap {
	f := make(template.FuncMap)

	/**
	* @param v float64 значение.
	 */
	f["displayfloat"] = func(v float64) string {
		return fmt.Sprintf("%6.2f", v)
	}

	return f
}
