package rss

import (
	"os"
)

var readFile = os.ReadFile
var stat = os.Stat
var exists = Exists

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
