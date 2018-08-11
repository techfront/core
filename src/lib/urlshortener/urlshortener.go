/**
* Библиотека для сокращения URL.
 */

package urlshortener

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

/**
* Ключ Google Api, задан в конфиге.
 */
var key string

type ShortMsg struct {
	Kind    string `json:"kind"`
	Id      string `json:"id"`
	LongUrl string `json:"longUrl"`
}

/**
* Инициализация и конфигурирование.
 */
func Setup(config map[string]string) error {
	key = config["google_api_key"]

	return nil
}

/**
* Сокращение URL, возращает url на goo.gl или ошибку.
 */
func Shorten(url string) (string, error) {
	var ShortMsg *ShortMsg

	request := gorequest.New()

	gUrl := "https://www.googleapis.com/urlshortener/v1/url?key=" + key

	res, _, err := request.Post(gUrl).
		Set("Accept", "application/json").
		Set("Content-Type", "application/json").
		Send(`{"longUrl":"` + url + `"}`).End()

	if err != nil {
		return "", fmt.Errorf("#error Request error: %v", err)
	}

	if res.Status == "200 OK" {
		if err := json.NewDecoder(res.Body).Decode(&ShortMsg); err != nil {
			return "", fmt.Errorf("#error Decode error: %v", err)
		}

		return ShortMsg.Id, nil
	}

	return "", fmt.Errorf("#error Request error: %v", res.Status)
}
