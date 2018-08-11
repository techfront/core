package twitter

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mrjones/oauth"
)

var key, secret, accessToken, accessTokenSecret string

/**
* Инициализация и конфигурирование.
 */
func Setup(k, s, at, ats string) error {
	if len(k) == 0 || len(s) == 0 || len(at) == 0 || len(ats) == 0 {
		return fmt.Errorf("#error setting secrets, null value")
	}

	key = k
	secret = s
	accessToken = at
	accessTokenSecret = ats

	return nil
}

// Tweet sends a status update to twitter - returns the response body or error
func Tweet(s string) ([]byte, error) {
	consumer := oauth.NewConsumer(key, secret, oauth.ServiceProvider{})
	//consumer.Debug(true)
	token := &oauth.AccessToken{Token: accessToken, Secret: accessTokenSecret}

	url := "https://api.twitter.com/1.1/statuses/update.json"
	data := map[string]string{"status": s}

	response, err := consumer.Post(url, data, token)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	fmt.Println("Response:", response.StatusCode, response.Status)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Check for unexpected status codes, and report them
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("#error sending tweet, unexpected status:%d\n\n%s\n", response.StatusCode, body)
	}

	return body, nil
}
