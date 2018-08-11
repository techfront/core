package view

import (
	"bytes"
	"fmt"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/kernel/view/helper"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	mu sync.RWMutex
	viewConfig View
	viewTemplates []string
	viewModifiers []ModifyFunc
	viewHelpers template.FuncMap
)

type ModifyFunc func(Context, *View)

type Context interface {
	Path() string
	Writer() http.ResponseWriter
	Message() router.Message
	Get(key string) interface{}
}

type View struct {
	Extension  string
	Layout     string
	Format     string
	Folder     string
	Production bool
	Vars       map[string]interface{}

	Modifiers []ModifyFunc
	Helpers   template.FuncMap

	Templates          []string
	templateCollection map[string]*template.Template

	muModifiers sync.RWMutex
	muTemplates sync.RWMutex
	muHelpers   sync.RWMutex
}

/**
* Функция New инициализирует View.
 */
func New(templateList ...string) *View {
	mu.Lock()
	defer mu.Unlock()

	v := &View{}
	v.Modifiers = viewConfig.Modifiers
	v.Helpers = viewConfig.Helpers
	v.Extension = viewConfig.Extension
	v.Layout = viewConfig.Layout
	v.Folder = viewConfig.Folder
	v.Format = viewConfig.Format
	v.Production = viewConfig.Production
	v.templateCollection = viewConfig.templateCollection

	v.Templates = append(viewTemplates, templateList...)
	v.Vars = make(map[string]interface{}, 0)

	return v
}

/**
* Настройка параметров.
 */
func Setup(config View) {
	viewConfig = config
	viewConfig.Modifiers = viewModifiers
	viewConfig.Helpers = viewHelpers
	viewConfig.Templates = []string{}
	viewConfig.templateCollection = make(map[string]*template.Template)
}

func LoadTemplates(temps []string) {
	mu.RLock()
	defer mu.RUnlock()

	viewTemplates = temps
}

func LoadModifiers(fn ...ModifyFunc) {
	mu.RLock()
	defer mu.RUnlock()

	viewModifiers = fn
}

func LoadHelpers(fms ...template.FuncMap) {
	mu.RLock()
	defer mu.RUnlock()

	fm := make(template.FuncMap)

	defaultHelpers := []template.FuncMap{
		helper.SanitizeHelper(),
		helper.DateHelper(),
		helper.XMLPreambleHelper(),
		helper.EscapeHelper(),
		helper.EscapeURLHelper(),
		helper.SetHelper(),
		helper.FieldHelper(),
		helper.TextareaHelper(),
		helper.SelectHelper(),
		helper.EmptyHelper(),
		helper.SafeHelper(),
		helper.AttrHelper(),
		helper.UrlHelper(),
	}

	fms = append(fms, defaultHelpers...)

	for _, m := range fms {
		for k, v := range m {
			fm[k] = v
		}
	}

	viewHelpers = fm
}

/**
* Компилирование шаблона
 */
func (v *View) Render(c Context) error {

	// Создание ключа для доступа к уже скомпилированным шаблонам
	key := strings.Join(v.Templates, ":")

	// Получение скомпилированного шаблона из коллекции
	v.muTemplates.RLock()
	tc, ok := v.templateCollection[key]
	v.muTemplates.RUnlock()

	// Получение хелперов
	v.muHelpers.RLock()
	hc := v.Helpers
	v.muHelpers.RUnlock()

	if !ok || !v.Production {
		var templateList []string
		templateList = append(templateList, v.Templates...)

		// Поиск шаблонов по заданным путям
		for i, name := range templateList {
			path, err := filepath.Abs(v.Folder + string(os.PathSeparator) + name + "." + v.Extension)
			if err != nil {
				return fmt.Errorf("#error No such template found %s", name+"."+v.Extension)
			}

			templateList[i] = path
		}

		templates, err := template.New(key).Funcs(hc).ParseFiles(templateList...)
		if err != nil {
			return fmt.Errorf("#error Template Parse Error: %s", err.Error())
		}

		// Добавляем в коллекцию
		v.muTemplates.Lock()
		v.templateCollection[key] = templates
		v.muTemplates.Unlock()

		tc = templates
	}

	// Получение модификаторов
	v.muModifiers.RLock()
	mc := v.Modifiers
	v.muModifiers.RUnlock()

	// Инициализация модификаторов
	for _, fn := range mc {
		fn(c, v)
	}

	layout := v.Layout + "." + v.Extension
	w := c.Writer()
	w.Header().Set("Content-Type", v.Format+"; charset=utf-8")

	err := tc.Funcs(hc).ExecuteTemplate(w, layout, v.Vars)
	if err != nil {
		return fmt.Errorf("#error Template File Error: %s", err.Error())
	}

	return nil
}

/**
* Функция RenderToString компилирует единственный шаблон в строку.
* Используется для рендеринга Email-шаблонов.
 */
func (v *View) RenderToString() (string, error) {

	var tmplPath string
	var tmplName string

	tmplPath = v.Templates[0]

	key := strings.Join(v.Templates, ":")

	// Получение хелперов
	v.muHelpers.RLock()
	hc := v.Helpers
	v.muHelpers.RUnlock()

	// Получение абсолютного пути
	path, err := filepath.Abs(v.Folder + string(os.PathSeparator) + tmplPath + "." + v.Extension)
	if err != nil {
		return "", fmt.Errorf("#error No such template found %s", tmplPath+"."+v.Extension)
	}
	tmplPath = path
	tmplName = filepath.Base(path)

	// Компиляция шаблона
	tc, err := template.New(key).Funcs(hc).ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("#error Template Parse Error: %s", err.Error())
	}

	// Получение названия шаблона
	tmpl := tc.Lookup(tmplName)
	if tmpl == nil {
		return "", fmt.Errorf("#error Loading template for %s", tmplName)
	}

	// Рендеринг
	var rendered bytes.Buffer

	err = tmpl.Funcs(hc).Execute(&rendered, v.Vars)
	if err != nil {
		return "", fmt.Errorf("#error Template File Error: %s", err.Error())
	}

	// Запись буфера в строку
	content := rendered.String()

	return content, nil
}
