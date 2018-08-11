package cleanhtml

import (
	"html/template"
	"github.com/kennygrant/sanitize"
)

var defaultTags = []string{"h1", "h2", "h3", "h4", "h5", "h6", "div", "span", "hr", "p", "br", "b", "i", "strong", "em", "ol", "ul", "li", "a", "img", "pre", "code", "blockquote"}
var defaultAttributes = []string{"id", "class", "src", "href", "title", "alt", "name", "rel", "height", "width", "srcset", "data-message", "data-modal", "data-close", "data-user-id"}

/**
 * Хелпер принимает строку и возвращает чистый HTML.
 */
func Helper() template.FuncMap {
	f := make(template.FuncMap)

	f["cleanhtml"] = func(s string) template.HTML {
		s, err := sanitize.HTMLAllowing(s, defaultTags, defaultAttributes)
		if err != nil {
			return template.HTML("")
		}

		return template.HTML(s)
	}

	return f
}
