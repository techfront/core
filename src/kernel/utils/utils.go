package utils

import (	
	"fmt"
)

/**
* Функция преобразует значение в строку.
*
* @param interface{}
* @return string
*/
func ToString(v interface{}) string {
	var s string

	switch v.(type) {
	case uint8, int8, uint16, int16, uint32, int32, uint64, int64, int, uint:
		s = fmt.Sprintf("%d", v)
	case float32, float64, complex64:
		s = fmt.Sprintf("%g", v)
	case bool:
		s = fmt.Sprintf("%t", v)
	case string:
		s = v.(string)
	case []byte:
		s = string(v.([]byte))
	default:
		s = ""
	}

	return s
}

/**
* Функция ToByte преобразует значение в байтмассив.
*
* @param interface{}
* @return string
*/
func ToByte(v interface{}) []byte {
	var b []byte

	switch v.(type) {
	case string:
		if IsSidHex(v.(string)) {
			v = SidHex(v.(string))
		}
		b = []byte(v.(string))
	case Sid:
		b = []byte(v.(Sid))
	case []byte:
		b = v.([]byte)
	default:
		b = nil
	}

	return b
}