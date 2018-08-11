package cache

import (
	"time"
)

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

/**
* Janitor выполняет функцию clearExpired() по интервалу.
*
* @param c *Cache
 */
func (j *janitor) Run(c *Cache) {
	j.stop = make(chan bool)
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			clearExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

func runJanitor(c *Cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
	}

	c.janitor = j

	go j.Run(c)
}
