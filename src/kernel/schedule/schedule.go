package schedule

import (
	"time"
)

type Context interface {
	Config(key string) string
	Production() bool
	Set(key string, data interface{})
	Get(key string) interface{}
	Log(message string)
	Logf(format string, v ...interface{})
}

type ScheduledFunc func(Context)

var logger Logger
var config Config

func Setup(l Logger, c Config) {
	logger = l
	config = c
}

func At(f ScheduledFunc, context Context, t time.Time, i time.Duration) chan struct{} {
	task := make(chan struct{})
	now := time.Now().UTC()

	for now.Sub(t) > 0 {
		t = t.Add(i)
	}

	if !context.Production() {
		context.Logf("schedule: task registered for:%s", t)
	}

	tillTime := t.Sub(now)
	time.AfterFunc(tillTime, func() {
		go f(context)

		if i > 0 {
			ticker := time.NewTicker(i)
			go func() {
				for {
					select {
					case <-ticker.C:
						go f(context)
					case <-task:
						ticker.Stop()
						return
					}
				}
			}()
		}
	})

	return task
}
