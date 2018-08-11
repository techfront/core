package gzip

import (
	"github.com/klauspost/pgzip"
	"io/ioutil"
	"sync"
)

const (
	BestCompression = pgzip.BestCompression
	BestSpeed = pgzip.BestSpeed
	DefaultCompression = pgzip.DefaultCompression
	NoCompression = pgzip.NoCompression
)

var Pool sync.Pool

func Setup() {
	Pool.New = func() interface{} {
		gz, err := pgzip.NewWriterLevel(ioutil.Discard, DefaultCompression)
		if err != nil {
			return err
		}

		return gz
	}
}