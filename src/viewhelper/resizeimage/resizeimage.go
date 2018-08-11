package resizeimage

import (
	"html/template"
)

func Helper() template.FuncMap {
	f := make(template.FuncMap)

	f["resizeimage"] = func(url string, size string) string {

		//newUrl := "https://techfront.org/api/imageproxy/" + size + "/uploads" + url

		newUrl := "/uploads" + url + "?size=" + size

		return newUrl
	}

	return f
}
