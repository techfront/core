/**
* Библиотека для работы с API Embedly.
 */

package embedly

import (
	"encoding/json"
	"fmt"
	"github.com/techfront/core/src/lib/upload"
	"net/http"
)

/**
* Переменная содержит ключ Embedly API.
 */
var key string

type EmbedlyResponse struct {
	ThumbnailURL    string `json:"thumbnail_url"`
	ThumbnailHeight int    `json:"thumbnail_height"`
	ThumbnailWidth  int    `json:"thumbnail_width"`
}

/**
 * Инициализация и конфигурирование.
 */
func Setup(config map[string]string) {
	key = config["embedly_key"]
}

/**
 * Функция получает Thumbnail по URL.
 */
func getEmbedlyThumbnail(url string) (string, error) {
	var EmbedlyResponse *EmbedlyResponse

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("#error request error: %v", err)
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&EmbedlyResponse)

	return EmbedlyResponse.ThumbnailURL, nil
}

/**
 * Функция загружает Thumbnail по URL.
 *
 * TODO: Сомневаюсь, что это здесь нужно.
 */
func UploadEmbedlyThumbnail(url string) (string, error) {
	api := "http://api.embedly.com/1/oembed?url=" + url + "&key=" + key
	res, err := getEmbedlyThumbnail(api)
	if err != nil {
		return "", err
	}

	path, err := upload.UploadFromUrl(res, "thumbnail")
	if err != nil {
		return "", fmt.Errorf("#error upload error: %v", err)
	}

	return path, nil
}
