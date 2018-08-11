package vk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

type saveWallPhoto struct {
	Response []struct {
		Aid      int    `json:"aid"`
		Created  int    `json:"created"`
		Height   int    `json:"height"`
		ID       string `json:"id"`
		OwnerID  int    `json:"owner_id"`
		Pid      int    `json:"pid"`
		Src      string `json:"src"`
		SrcBig   string `json:"src_big"`
		SrcSmall string `json:"src_small"`
		Text     string `json:"text"`
		Width    int    `json:"width"`
	} `json:"response"`
}

type getWallUploadServer struct {
	Response struct {
		AlbumId       int    `json:"album_id"`
		UserId      int    `json:"user_id"`
		UploadURL string `json:"upload_url"`
	} `json:"response"`
}

type uploadResponse struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

var accessToken, groupId string

// Setup sets our secret keys
func Setup(at, gi string) error {
	if len(at) == 0 || len(gi) == 0 {
		return fmt.Errorf("#error setting secrets, null value")
	}

	accessToken = at
	groupId = gi

	return nil
}

func photoUpload(url_source string) (string, error) {
	var b bytes.Buffer

	var uploadServer *getWallUploadServer
	var uploadResponse *uploadResponse
	var saveWallPhoto *saveWallPhoto

	// Get upload server url
	urlUploadServer := "https://api.vk.com/method/photos.getWallUploadServer?v=5.73&group_id=" + groupId + "&access_token=" + accessToken
	resUploadServer, _ := http.Get(urlUploadServer)
	bytesUploadServer, err := ioutil.ReadAll(resUploadServer.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(bytesUploadServer), &uploadServer); err != nil {
		return "", err
	}
	defer resUploadServer.Body.Close()

	// Upload Photo from url
	resSourceImg, err := http.Get(url_source)
	if err != nil {
		return "", err
	}
	defer resSourceImg.Body.Close()

	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("photo", "photo.jpg")
	if err != nil {
		return "", fmt.Errorf("CreateFormFile error: %v", err)
	}

	if _, err = io.Copy(fw, resSourceImg.Body); err != nil {
		return "", fmt.Errorf("Copy error: %v", err)
	}
	w.Close()

	reqPostImg, err := http.NewRequest("POST", uploadServer.Response.UploadURL, &b)
	if err != nil {
		return "", fmt.Errorf("NewRequest error: %v", err)
	}

	reqPostImg.Header.Set("Content-Type", w.FormDataContentType())
	clientPostImg := &http.Client{}
	resPostImg, errPostImg := clientPostImg.Do(reqPostImg)
	if errPostImg != nil {
		return "", fmt.Errorf("Do error: %v", errPostImg)
	}
	defer resPostImg.Body.Close()

	if err := json.NewDecoder(resPostImg.Body).Decode(&uploadResponse); err != nil {
		return "", fmt.Errorf("Decode error: %v", err)
	}

	// Save Wall Photo
	data := url.Values{}
	data.Set("server", strconv.Itoa(uploadResponse.Server))
	data.Add("photo", uploadResponse.Photo)
	data.Add("hash", uploadResponse.Hash)
	data.Add("group_id", groupId)

	urlSaveImg := "https://api.vk.com/method/photos.saveWallPhoto?v=5.73&access_token=" + accessToken

	resSaveImg, err := http.Post(urlSaveImg, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("Save error: %v", err)
	}
	defer resSaveImg.Body.Close()

	if err := json.NewDecoder(resSaveImg.Body).Decode(&saveWallPhoto); err != nil {
		return "", fmt.Errorf("Decode error: %v", err)
	}

	// Return Photo Id
	return saveWallPhoto.Response[0].ID, nil
}

func Post(message string, thumbnail string, link string) error {
	var attachments string

	attachments = link

	if len(thumbnail) > 0 {
		photoId, err := photoUpload(thumbnail)
		if err != nil {
			return err
		}

		attachments = photoId

		if len(link) > 0 {
			attachments = photoId + "," + link
		}
	}

	data := url.Values{}
	data.Set("owner_id", "-" + groupId)
	data.Set("from_group", "1")
	data.Add("message", message)
	data.Add("attachments", attachments)

	res, err := http.PostForm("https://api.vk.com/method/wall.post?v=5.73&access_token=" + accessToken, data)
	if err != nil {
		return fmt.Errorf("Post error: %v", err)
	}

	defer res.Body.Close()

	return nil
}
