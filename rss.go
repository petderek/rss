package rss

import (
	"encoding/xml"
)

type image struct {
	Url   string `xml:"url"`
	Link  string `xml:"link"`
	Title string `xml:"title"`
}

type item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PublishDate string `xml:"pubDate"`
	Link        string `xml:"link"`
	Content     string `xml:"content"`
	Guid        string `xml:"guid"`
}

type channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Image       image  `xml:"image"`
	Items       []item `xml:"item"`
}

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel channel  `xml:"channel"`
}

func parseXML(data []byte) (rss, error) {
	v := rss{}
	err := xml.Unmarshal(data, &v)
	return v, err
}
