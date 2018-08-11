package validate

import (
	"fmt"
	"strconv"
	"time"
)

func Float(param interface{}) float64 {
	var v float64
	if param != nil {
		switch param.(type) {
		case int64:
			v = float64(param.(int64))
		default:
			v = param.(float64)
		}
	}

	return v
}

func Boolean(param interface{}) bool {
	var v bool
	if param != nil {
		v = param.(bool)
	}

	return v
}

func Int(param interface{}) int64 {
	var v int64
	if param != nil {
		switch param.(type) {
		case float64:
			v = int64(param.(float64))
		case int:
			v = int64(param.(int))
		default:
			v = param.(int64)
		}
	}

	return v
}

func String(param interface{}) string {
	var v string
	if param != nil {
		v = param.(string)
	}
	return v
}

func Time(param interface{}) time.Time {
	var v time.Time
	if param != nil {
		v = param.(time.Time)
	}

	return v
}

func Length(param string, min int, max int) error {
	length := len(param)
	if min != -1 && length < min {
		return fmt.Errorf("length of string %s %d, expected > %d", param, length, min)
	}
	if max != -1 && length > max {
		return fmt.Errorf("length of string %s %d, expected < %d", param, length, max)
	}

	return nil
}

func Within(param string, min float64, max float64) error {
	f, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return fmt.Errorf("invalid float param %s", param)
	}
	if f < min {
		return fmt.Errorf("%0.2f is less than minimum %0.2f", f, min)
	}
	if f > max {
		return fmt.Errorf("%0.2f is more than maximum %0.2f", f, max)
	}

	return nil
}