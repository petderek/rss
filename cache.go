package rss

import (
	"errors"
	"path/filepath"
)

var (
	ErrNotFound = errors.New("not found")
)

type Cache struct {
	dir string
}

func NewCache(dir string) *Cache {
	return &Cache{dir: dir}
}

func (c *Cache) Get(subscription string) (*Subscription, error) {
	target := filepath.Join(c.dir, subscription)
	if !exists(target) {
		return nil, ErrNotFound
	}
	return &Subscription{dir: target}, nil
}

type Subscription struct {
	name string
	dir  string
}

func (s *Subscription) GetRssData() ([]byte, error) {
	return readFile(filepath.Join(s.dir, s.name, "rss.xml"))
}

func (s *Subscription) GetGuidContent(guid string) ([]byte, error) {
	return nil, nil
}
