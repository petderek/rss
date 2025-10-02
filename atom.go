package rss

import (
	"encoding/xml"
)

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type AtomContent struct {
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}

type AtomEntry struct {
	ID        string      `xml:"id"`
	Title     string      `xml:"title"`
	Summary   string      `xml:"summary"`
	Content   AtomContent `xml:"content"`
	Link      []AtomLink  `xml:"link"`
	Published string      `xml:"published"`
	Updated   string      `xml:"updated"`
}

type AtomFeed struct {
	XMLName xml.Name     `xml:"feed"`
	Title   string       `xml:"title"`
	Subtitle string      `xml:"subtitle"`
	Link    []AtomLink   `xml:"link"`
	Icon    string       `xml:"icon"`
	Entries []AtomEntry  `xml:"entry"`
}

func parseAtomXML(data []byte) (AtomFeed, error) {
	v := AtomFeed{}
	err := xml.Unmarshal(data, &v)
	return v, err
}