package router

import (
	"path"
	"path/filepath"
)

/**
* Функция приводит путь к приемлимому виду.
 */
func canonicalPath(p string) string {
	canonical := path.Clean(p)
	if len(canonical) == 0 {
		canonical = "/"
	} else if canonical[0] != '/' {
		canonical = "/" + canonical
	}

	return canonical
}

/**
 * Функция pathExt получает расширение ресурса.
 */
func pathExt(p string) string {
	return filepath.Ext(p)
}

/**
* Функция возвращает True, если значение(int64) есть в массиве.
 */
func contains(list []int64, i int64) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}
	return false
}

/**
* Функция возвращает True, если строка есть в массиве.
 */
func containsString(allowed []string, s string) bool {
	for _, v := range allowed {
		if v == s {
			return true
		}
	}
	return false
}
