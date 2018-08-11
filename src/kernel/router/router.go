package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

/**
* Интерфейс Handler необходим для обработки HTTP запросов.
 */
type Handler interface {
	ServeHTTP(Context) error
}

/**
* Тип HandlerFunc определяет функцию для обработки HTTP запросов.
 */
type HandlerFunc func(Context) error

/**
* ServeHTTP реализует интерфейс для обработки HTTP запросов.
 */
func (f HandlerFunc) ServeHTTP(c Context) error {
	return f(c)
}

/**
* Тип ErrorHandler определяет функцию для обработки ошибок.
 */
type ErrorHandler func(Context, error)

/**
* Интерфейс Config необходим, что бы получить конфигурацию веб-сервера.
 */
type Config interface {
	Production() bool
	Config(string) string
}

/**
* Интерфейс Logger необходим, что бы журналировать HTTP Запросы.
 */
type Logger interface {
	Printf(format string, args ...interface{})
}

type Router struct {
	// Защита для посредников и запросов
	mu sync.RWMutex

	// HTTPRouter
	Router *httprouter.Router

	// Хендлер файлов
	FileHandler HandlerFunc

	// Хендлер ошибок
	ErrorHandler ErrorHandler

	// Конфигурация веб-сервера
	Config Config

	// Журналирование
	Logger Logger

	// Посредники
	middlewares []func(Handler) Handler

	// Пул для контекста
	pool sync.Pool
}

/**
* Функция New инициализирует роутер.
 */
func New(l Logger, s Config) *Router {

	r := &Router{
		Router:       httprouter.New(),
		FileHandler:  fileHandler,
		ErrorHandler: errorHandler,
		Config:       s,
		Logger:       l,
	}

	r.pool.New = func() interface{} {
		return r.buildContext(nil, nil, nil)
	}

	// Переопределение обработчика ошибок MethodNotAllowed для httprouter
	r.Router.MethodNotAllowed = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		// Построение контекста
		context := r.pool.Get().(*ConcreteContext)
		defer r.pool.Put(context)
		context.Reset(writer, request, nil)

		// Запуска обработчика ошибок
		r.ErrorHandler(context, MethodNotAllowedError(nil))
	})

	// Переопределение обработчика ошибок NotFound для httprouter
	r.Router.NotFound = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		// Построение контекста
		context := r.pool.Get().(*ConcreteContext)
		defer r.pool.Put(context)
		context.Reset(writer, request, nil)

		// Вначале, поиск по файлам, если файла нет, то обрабатываем как ошибку.
		err := r.FileHandler(context)
		if err != nil {
			r.ErrorHandler(context, err)
		}
	})

	http.Handle("/", r.Router)

	return r
}

/**
* Функция Logf отправляет сообщение с аргументами.
 */
func (r *Router) Logf(format string, v ...interface{}) {
	r.Logger.Printf(format, v...)
}

/**
* Функция Log отправляет сообщение.
 */
func (r *Router) Log(message string) {
	r.Logf(message)
}

func (r *Router) Delete(path string, h HandlerFunc) {
	r.Router.DELETE(path, r.Handle(h))
}

func (r *Router) Get(path string, h HandlerFunc) {
	r.Router.GET(path, r.Handle(h))
}

func (r *Router) Head(path string, h HandlerFunc) {
	r.Router.HEAD(path, r.Handle(h))
}

func (r *Router) Options(path string, h HandlerFunc) {
	r.Router.OPTIONS(path, r.Handle(h))
}

func (r *Router) Patch(path string, h HandlerFunc) {
	r.Router.PATCH(path, r.Handle(h))
}

func (r *Router) Post(path string, h HandlerFunc) {
	r.Router.POST(path, r.Handle(h))
}

func (r *Router) Put(path string, h HandlerFunc) {
	r.Router.PUT(path, r.Handle(h))
}

/**
* Функция AddMiddleware выполняет добавление нового посредника.
 */
func (r *Router) AddMiddleware(m func(Handler) Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Добавляем посредник в массив
	r.middlewares = append(r.middlewares, m)
}

/**
* Функция buildChain строит цепочку посредников.
 */
func (r *Router) buildChain(h Handler) Handler {
	r.mu.Lock()
	defer r.mu.Unlock()

	var chain Handler

	chain = h

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		chain = r.middlewares[i](chain)
	}

	return chain
}

/**
* Функция buildContext выполняет построение контекста.
 */
func (r *Router) buildContext(writer http.ResponseWriter, request *http.Request, p httprouter.Params) Context {
	return &ConcreteContext{
		request: request,
		writer:  writer,
		logger:  r.Logger,
		config:  r.Config,
		params:  p,
		path:    "",
		ext:     "",
		data:    make(map[string]interface{}, 0),
		message: Message{},
	}
}

/**
 * Функция Handle это обёртка для интеграции с httprouter.
 */
func (r *Router) Handle(h Handler) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {

		// Построение контекста
		context := r.pool.Get().(*ConcreteContext)
		defer r.pool.Put(context)
		context.Reset(writer, request, p)

		// Построение цепочки посредников.
		chain := r.buildChain(h)

		if err := chain.ServeHTTP(context); err != nil {
			r.ErrorHandler(context, err)
			return
		}
	}
}
