package rss

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	readFile  = os.ReadFile
	exists    = Exists
	fetch     = Fetch
	mkdir     = os.MkdirAll
	writeFile = os.WriteFile
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Fetch(url *url.URL) ([]byte, error) {
	if url == nil {
		return nil, errors.New("url is nil")
	}
	r, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
