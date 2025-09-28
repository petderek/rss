package rss

import (
	"os"
)

type deps struct {
	ReadFile func(string) ([]byte, error)
	Stat     func(string) (os.FileInfo, error)
	Exists   func(string) bool
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func defaultDeps() deps {
	return deps{
		ReadFile: os.ReadFile,
		Stat:     os.Stat,
		Exists:   Exists,
	}
}
