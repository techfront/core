package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Context interface {
	Request() *http.Request

	Writer() http.ResponseWriter

	SetRequest(*http.Request)

	SetWriter(http.ResponseWriter)

	ClientIP() string

	Path() string

	Ext() string

	Config(key string) string

	Production() bool

	Params() (Params, error)

	Param(key string) string

	Data() map[string]interface{}

	Set(key string, data interface{})

	Get(key string) interface{}

	SetMessage(uid string, body string, style int64)

	Message() Message

	Log(message string)

	Logf(format string, v ...interface{})

	Reset(writer http.ResponseWriter, request *http.Request, p httprouter.Params)
}

type ConcreteContext struct {
	request *http.Request

	writer http.ResponseWriter

	path string

	ext string

	logger Logger

	config Config

	params httprouter.Params

	data map[string]interface{}

	message Message
}

func (c *ConcreteContext) Logf(format string, v ...interface{}) {
	c.logger.Printf(format, v...)
}

func (c *ConcreteContext) Log(message string) {
	c.Logf(message)
}

func (c *ConcreteContext) Request() *http.Request {
	return c.request
}

func (c *ConcreteContext) Writer() http.ResponseWriter {
	return c.writer
}

func (c *ConcreteContext) SetRequest(r *http.Request) {
	c.request = r
}

func (c *ConcreteContext) SetWriter(w http.ResponseWriter) {
	c.writer = w
}

func (c *ConcreteContext) Path() string {
	return c.path
}

func (c *ConcreteContext) Ext() string {
	return c.ext
}

func (c *ConcreteContext) Config(key string) string {
	return c.config.Config(key)
}

func (c *ConcreteContext) Production() bool {
	return c.config.Production()
}

func (c *ConcreteContext) Data() map[string]interface{} {
	return c.data
}

func (c *ConcreteContext) Set(key string, data interface{}) {
	c.data[key] = data
}

func (c *ConcreteContext) Get(key string) interface{} {
	return c.data[key]
}

/**
* Функция SetMessage добавляет сообщение в контекст.
*
* @param uid string уникальный идентификатор.
* @param body string тело сообщения.
* @param style int64 тип сообщения.
 */
func (c *ConcreteContext) SetMessage(uid string, body string, style int64) {
	c.message = Message{
		uid: uid,
		style: style,
		body:  body,
	}
}

/**
* Функция Message получаем сообщение из контекста.
*
* @return Message сообщение.
 */
func (c *ConcreteContext) Message() Message {
	return c.message
}

/**
* Функция parseRequest получает параметры отправленные из формы для текущего запроса.
 */
func (c *ConcreteContext) parseRequest() error {

	if c.request.Body == nil {
		return nil
	}

	err := c.request.ParseForm()
	if err != nil {
		return err
	}

	return nil
}

/**
* Функция Params параметр по ключу.
 */
func (c *ConcreteContext) Param(key string) string {
	params, err := c.Params()
	if err != nil {
		c.Logf("Error parsing request %s", err)
		return ""
	}

	return params.Get(key)
}

/**
* Функция Params получает все параметры текущего запроса.
 */
func (c *ConcreteContext) Params() (Params, error) {
	params := Params{}

	if c.request.Form == nil {
		err := c.parseRequest()
		if err != nil {
			c.Logf("Error parsing request params %s", err)
			return nil, err
		}

	}

	for k, v := range c.request.Form {
		for _, vv := range v {
			params.Add(k, vv)
		}
	}

	routeParams := c.params
	for i := range routeParams {
		params.Add(routeParams[i].Key, routeParams[i].Value)
	}

	return params, nil
}

/**
* Функция ClientIP возвращает реальный IP клиента.
*
* @return string
 */
func (c *ConcreteContext) ClientIP() string {
	address := c.request.Header.Get("X-Real-IP")
	if len(address) > 0 {
		return address
	}

	address = c.request.Header.Get("X-Forwarded-For")
	if len(address) > 0 {
		return address
	}

	return c.request.RemoteAddr
}

/**
* Функция Reset конфигурирует (реконфигурирует) контекст.
*
* @param writer http.ResponseWriter
* @param *http.Request
* @param p httprouter.Params
 */
func (c *ConcreteContext) Reset(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {
	c.request = request
	c.writer = writer
	c.params = p
	c.path = canonicalPath(request.URL.Path)
	c.ext = pathExt(request.URL.Path)
	c.data = make(map[string]interface{}, 0)
	c.message = Message{}
}
