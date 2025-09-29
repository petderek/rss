package rss

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

const sample = `<rss>
<channel>
<title>title</title>
<description>desc</description>
<link>http://example.com/title</link>
<image>
<link>http://example.com/imagelink</link>
<url>http://example.com/imageurl</url>
<title>imagetitle</title>
</image>
<item>
<title>item1</title>
<description>item1desc</description>
<link>http://example.com/item1link</link>
<content:encoded>item1content</content:encoded>
<pubDate>Fri, 19 Sep 2025 12:05:57</pubDate>
<guid>item1guid</guid>
</item>
<item>
<title>item2</title>
<description>item2desc</description>
<link>http://example.com/item2link</link>
<content:encoded>item2content</content:encoded>
<pubDate>Fri, 19 Sep 2025 12:05:57</pubDate>
<guid>item2guid</guid>
</item>
</channel>
</rss>`

var (
	expected = Rss{
		XMLName: xml.Name{Local: "rss"},
		Channel: Channel{
			Title:       "title",
			Description: "desc",
			Link:        "http://example.com/title",
			Image: Image{
				Url:   "http://example.com/imageurl",
				Link:  "http://example.com/imagelink",
				Title: "imagetitle",
			},
			Items: []Item{
				{
					Title:       "item1",
					Description: "item1desc",
					Link:        "http://example.com/item1link",
					Content:     "item1content",
					Guid:        "item1guid",
					PublishDate: "Fri, 19 Sep 2025 12:05:57",
				},
				{
					Title:       "item2",
					Description: "item2desc",
					Link:        "http://example.com/item2link",
					Content:     "item2content",
					Guid:        "item2guid",
					PublishDate: "Fri, 19 Sep 2025 12:05:57",
				},
			},
		},
	}
)

func TestParseXML(t *testing.T) {
	c, err := parseXML([]byte(sample))
	require.NoError(t, err)
	assert.Equal(t, expected, c)
}

func parseUrl(s string) url.URL {
	u, _ := url.Parse(s)
	return *u
}
