package helper

import (
	"fmt"
	"html/template"
)

// TODO: Придумать реализацию с интерфейсом.

/* type ConcreteSelect struct {
	Name string
	Label string
	Value int64
	Options []ConcreteOption
} */

/**
* Тип Option реализует интерфейс для поля Option.
 */

/* type Option interface {
	GetId() int64
	GetName() string
	SetId(int64)
	SetName(string)
} */

/**
* Тип ContreteOption это структура конкретного Option.
 */
type ConcreteOption struct {
	Id   int64
	Name string
}

/**
* GetId возвращает ID для текущего Option.
 */
func (o *ConcreteOption) GetId() int64 {
	return o.Id
}

/**
* GetName возвращает Name для текущего Option.
 */
func (o *ConcreteOption) GetName() string {
	return o.Name
}

/**
* SetId задает Id для текущего Option.
 */
func (o *ConcreteOption) SetId(id int64) {
	o.Id = id
}

/**
* SetName задает Name для текущего Option.
 */
func (o *ConcreteOption) SetName(name string) {
	o.Name = name
}

/**
* SelectHelper создаёт поле вида Select в шаблоне.
 */
func SelectHelper() template.FuncMap {
	f := make(template.FuncMap)

	f["select"] = func(label string, name string, value int64, options []ConcreteOption) template.HTML {
		tmpl := `<div class="field"><label>%s</label><select class="select" type="select" name="%s">%s</select></div>`

		if label == "" {
			tmpl = `%s<select type="select" name="%s">%s</select>`
		}

		opts := ""
		for _, o := range options {

			s := ""
			if o.GetId() == value {
				s = "selected"
			}

			opts += fmt.Sprintf(`<option value="%d" %s>%s</option>`, o.GetId(), s, template.HTMLEscapeString(o.GetName()))
		}

		output := fmt.Sprintf(tmpl, template.HTMLEscapeString(label), template.HTMLEscapeString(name), opts)

		return template.HTML(output)
	}

	return f
}
