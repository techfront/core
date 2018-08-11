package schedule

type Logger interface {
	Printf(format string, args ...interface{})
}

type Config interface {
	Production() bool
	Config(string) string
}

type ConcreteContext struct {
	logger Logger
	config Config
	data map[string]interface{}
}

func NewContext() *ConcreteContext {
	return &ConcreteContext{
		logger: logger,
		config: config,
		data: make(map[string]interface{}),
	}
}

func (c *ConcreteContext) Logf(format string, v ...interface{}) {
	c.logger.Printf(format, v...)
}

func (c *ConcreteContext) Log(message string) {
	c.Logf(message)
}

func (c *ConcreteContext) Config(key string) string {
	return c.config.Config(key)
}

func (c *ConcreteContext) Production() bool {
	return c.config.Production()
}

func (c *ConcreteContext) Set(key string, data interface{}) {
	c.data[key] = data
}

func (c *ConcreteContext) Get(key string) interface{} {
	return c.data[key]
}