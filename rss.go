package fuse_rss

import (
	"net/url"
	"time"
)

type image struct {
	url   url.URL
	link  url.URL
	title string
}

type item struct {
	title       string
	description string
	published   time.Time
	link        url.URL
	content     string
	guid        string
}

type channel struct {
	title       string
	link        string
	description string
	image       image
	items       []item
}
