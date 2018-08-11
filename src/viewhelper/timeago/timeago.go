package timeago

import (
	"fmt"
	"html/template"
	"time"
)

/**
* Дата и время основания проекта (1 января 1970).
 */
var foundingTime = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

/**
* Функция RusMonth возвращает название месяца на русском языке.
*
* @param m time.Month
* @return string название месяца на русском языке.
 */
func RusMonth(m time.Month) (month string) {
	switch m {
	case time.January:
		month = "Января"
	case time.February:
		month = "Февраля"
	case time.March:
		month = "Марта"
	case time.April:
		month = "Апреля"
	case time.May:
		month = "Мая"
	case time.June:
		month = "Июня"
	case time.July:
		month = "Июля"
	case time.August:
		month = "Августа"
	case time.September:
		month = "Сентября"
	case time.October:
		month = "Октября"
	case time.November:
		month = "Ноября"
	case time.December:
		month = "Декабря"
	}

	return
}

/**
* Хелпер конвертирует время в удобный формат.
 */
func Helper() template.FuncMap {
	f := make(template.FuncMap)

	f["timeago"] = func(d time.Time) string {

		// Если время оснавания больше, то приравниваем.
		if foundingTime.Unix() >= d.Unix() {
			d = foundingTime
		}

		duration := time.Since(d)
		hours := duration / time.Hour

		switch {
		case duration < time.Minute:
			return fmt.Sprintf("только что")
		case duration < time.Hour:
			unit := "минут"
			minutes := duration / time.Minute
			if (minutes%10 == 1) && (minutes%100 != 11) {
				unit = "минуту"
			} else if (minutes%10 >= 2) && (minutes%10 <= 4) && (minutes%100 < 10 || minutes%100 >= 20) {
				unit = "минуты"
			}
			return fmt.Sprintf("%d %s назад", minutes, unit)

		case duration < time.Hour*24:
			unit := "часов"
			if (hours%10 == 1) && (hours%100 != 11) {
				unit = "час"
			} else if (hours%10 >= 2) && (hours%10 <= 4) && (hours%100 < 10 || hours%100 >= 20) {
				unit = "часа"
			}
			return fmt.Sprintf("%d %s назад", hours, unit)
		case duration < time.Hour*24*7:
			days := hours / 24
			unit := "дней"

			if (days%10 == 1) && (days%100 != 11) {
				unit = "день"
			} else if (days%10 >= 2) && (days%10 <= 4) && (days%100 < 10 || days%100 >= 20) {
				unit = "дня"
			}

			return fmt.Sprintf("%d %s назад", days, unit)
		case duration < time.Hour*24*30:
			unit := "недель"
			weeks := hours / 24 / 7

			if (weeks%10 == 1) && (weeks%100 != 11) {
				unit = "неделю"
			} else if (weeks%10 >= 2) && (weeks%10 <= 4) && (weeks%100 < 10 || weeks%100 >= 20) {
				unit = "недели"
			}

			return fmt.Sprintf("%d %s назад", weeks, unit)
		case duration < time.Hour*24*365:
			unit := "месяцев"
			months := hours / 24 / 30

			if (months%10 == 1) && (months%100 != 11) {
				unit = "месяц"
			} else if (months%10 >= 2) && (months%10 <= 4) && (months%100 < 10 || months%100 >= 20) {
				unit = "месяца"
			}

			return fmt.Sprintf("%d %s назад", months, unit)
		default:
			month := RusMonth(d.Month())
			layout := d.Format("02 %s 2006 в 15:04")

			return fmt.Sprintf(layout, month)
		}

	}

	return f
}
