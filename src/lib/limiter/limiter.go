/**
* Пакет limiter необходим для ограничения колличества запросов к ресурсам или действиям.
 */
package limiter

import (
	"gopkg.in/throttled/throttled.v2"
	"log"
)

var limiterStore *Store

type Limiter struct {
	gcra *throttled.GCRARateLimiter
}

/**
* Конфигурирование.
*
* @param config map[string]string
 */
func Setup(config map[string]string) {

	// Инициализация хранилища
	s, err := NewStore(config["limiter_store"])
	if err != nil {
		log.Fatal(err)
	}

	limiterStore = s
}

/**
* Функция New инициализирует Limiter.
*
* @param rate int кол-во запросов.
* @param measurement string величина времени.
* @param burst int всплеск (дополнительные запросы).
*
 */
func New(rate int, measurement string, burst int) (*Limiter, error) {

	var r throttled.Rate

	switch measurement {
	case "second":
		r = throttled.PerSec(rate)
	case "minute":
		r = throttled.PerMin(rate)
	case "hour":
		r = throttled.PerHour(rate)
	case "day":
		r = throttled.PerDay(rate)
	default:
		r = throttled.PerSec(rate)
	}

	quota := throttled.RateQuota{r, burst}
	gcra, err := throttled.NewGCRARateLimiter(limiterStore, quota)
	if err != nil {
		return nil, err
	}

	l := &Limiter{
		gcra: gcra,
	}

	return l, nil
}

func (l *Limiter) GCRA() *throttled.GCRARateLimiter {
	return l.gcra
}

func (l *Limiter) ByParams(p []string) throttled.HTTPRateLimiter {
	return throttled.HTTPRateLimiter{
		RateLimiter: l.gcra,
		VaryBy:      &throttled.VaryBy{Params: p},
	}
}

func (l *Limiter) ByPath() throttled.HTTPRateLimiter {
	return throttled.HTTPRateLimiter{
		RateLimiter: l.gcra,
		VaryBy:      &throttled.VaryBy{Path: true},
	}
}
