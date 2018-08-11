package assets

import (
	"fmt"
	"html/template"
	"strings"
)

const (
	styleTemplate  = `<link href="/assets/styles/%s" media="all" rel="stylesheet" type="text/css" />`
	scriptTemplate = `<script src="/assets/scripts/%s" type="text/javascript" ></script>`
)

// StyleLink converts a set of group names to one style link tag (production) or to a list of style link tags (development)
func (c *Collection) StyleLink() template.FuncMap {
	f := make(template.FuncMap)

	f["style"] = func(names ...string) template.HTML {
		var html template.HTML

		hash := c.Group("app").stylehash

		for _, name := range names {
			g := c.Group(name)
			if g.stylehash != "" {
				if c.serveCompiled {
					html = html + StyleLink(g.StyleName())
				} else {
					for _, f := range g.Styles() {
						html = html + StyleLink(f.name) + template.HTML("\n")
					}
				}
			} else {
				html = html + StyleLinkWithHash(name, hash)
			}
		}

		return html
	}

	return f
}

// ScriptLink converts a set of group names to one script tag (production) or to a list of script tags (development)
func (c *Collection) ScriptLink() template.FuncMap {
	f := make(template.FuncMap)

	f["script"] = func(names ...string) template.HTML {
		var html template.HTML

		hash := c.Group("app").scripthash

		for _, name := range names {
			g := c.Group(name)
			if g.scripthash != "" {
				if c.serveCompiled {
					html = html + ScriptLink(g.ScriptName())
				} else {
					for _, f := range g.Scripts() {
						html = html + ScriptLink(f.name) + template.HTML("\n")
					}
				}
			} else {
				html = html + ScriptLinkWithHash(name, hash)
			}
		}

		return html
	}

	return f
}

// StyleLinkWithHash returns an html tag for a given file path
func StyleLinkWithHash(name string, hash string) template.HTML {
	hash = template.URLQueryEscaper(hash)

	if !strings.HasSuffix(name, ".css") {
		name = fmt.Sprintf("%s.css?hash=%s", name, hash)
	} else {
		name = fmt.Sprintf("%s?hash=%s", name, hash)
	}

	return template.HTML(fmt.Sprintf(styleTemplate, name))
}

// ScriptLinkWithHash returns an html tag for a given file path
func ScriptLinkWithHash(name string, hash string) template.HTML {
	hash = template.URLQueryEscaper(hash)

	if !strings.HasSuffix(name, ".js") {
		name = fmt.Sprintf("%s.js?hash=%s", name, hash)
	} else {
		name = fmt.Sprintf("%s?hash=%s", name, hash)
	}

	return template.HTML(fmt.Sprintf(scriptTemplate, name))
}

// StyleLink returns an html tag for a given file path
func StyleLink(name string) template.HTML {
	if !strings.HasSuffix(name, ".css") {
		name = name + ".css"
	}

	return template.HTML(fmt.Sprintf(styleTemplate, template.URLQueryEscaper(name)))
}

// ScriptLink returns an html tag for a given file path
func ScriptLink(name string) template.HTML {
	if !strings.HasSuffix(name, ".js") {
		name = name + ".js"
	}

	return template.HTML(fmt.Sprintf(scriptTemplate, template.URLQueryEscaper(name)))
}
