package rss

import (
	"encoding/xml"
)

type Image struct {
	Url   string `xml:"url"`
	Link  string `xml:"link"`
	Title string `xml:"title"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PublishDate string `xml:"pubDate"`
	Link        string `xml:"link"`
	Content     string `xml:"encoded"`
	Guid        string `xml:"guid"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Image       Image  `xml:"image"`
	Items       []Item `xml:"item"`
}

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

func parseXML(data []byte) (Rss, error) {
	v := Rss{}
	err := xml.Unmarshal(data, &v)
	return v, err
}
