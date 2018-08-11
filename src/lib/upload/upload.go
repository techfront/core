package upload

import (
	"os"
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	"strconv"
	"crypto/md5"
	"encoding/hex"
	"github.com/h2non/filetype"
	"github.com/techfront/core/src/lib/resizer"
)

var uploadsDir string

func Setup(config map[string]string) {
	uploadsDir = config["uploads_dir"]
}

func getFileName() string {
	t := time.Now().UTC()
	s := fmt.Sprint(t)
	h := md5.Sum([]byte(s))

	return hex.EncodeToString(h[:])
}

func getFilePath() string {
	y, m, d := time.Now().Date()
	p := "/" + strconv.Itoa(y) + "/" + strconv.Itoa(int(m)) + "/" + strconv.Itoa(d) + "/"

	if dirExist := checkPathExist(uploadsDir + p); dirExist != true {
		os.MkdirAll(uploadsDir+p, 0755)
	}

	return p
}

func getPath(e string) string {
	path := getFilePath() + getFileName() + "." + e

	return path
}

func checkPathExist(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}

	return true
}

func deleteFile(p string) error {
	err := os.Remove(p)
	if err != nil {
		return err
	}

	return nil
}

func UploadFromUrl(url string, format string) (string, error) {

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	dump, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	fileType, unknown := filetype.Match(dump)
	if unknown != nil || !filetype.IsImage(dump) {
		return "", fmt.Errorf("#error File type is not valid")
	}

	extName := fileType.Extension
	path := getPath(extName)
	fullPath := uploadsDir + path

	if fileExist := checkPathExist(fullPath); fileExist != true {
		file, err := os.Create(fullPath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		_, err = file.Write(dump)
		if err != nil {
			return "", err
		}
	}

	if extName == "gif" {
		return path, nil
	}

	p, err := resizer.Resize(path, format)
	if err != nil {
		return "", err
	}

	if err := deleteFile(fullPath); err != nil {
		return "", err
	}

	return p, nil
}
