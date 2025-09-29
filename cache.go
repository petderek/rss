package rss

import (
	"errors"
	"path/filepath"
)

var (
	ErrNotFound = errors.New("not found")
)

const (
	rssxml = "rss.xml"
)

type Cache struct {
	dir string
}

func NewCache(dir string) *Cache {
	return &Cache{dir: dir}
}

func (c *Cache) GetSubscription(name string) (*FileEntry, error) {
	target := filepath.Join(c.dir, name, rssxml)
	if !exists(target) {
		return nil, ErrNotFound
	}
	return &FileEntry{filename: target}, nil
}

func (c *Cache) PutSubscription(name string, data []byte) error {
	target := filepath.Join(c.dir, name)
	if err := mkdir(target, 0755); err != nil {
		return err
	}
	target = filepath.Join(target, rssxml)
	return writeFile(target, data, 0644)
}

type FileEntry struct {
	filename string
}

func (s *FileEntry) Data() ([]byte, error) {
	return readFile(s.filename)
}
