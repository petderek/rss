package rss

func ToInternal(data []byte) (*Node, error) {
	r, err := parseXML(data)
	if err != nil {
		return nil, err
	}
	return fromRss(r), nil
}
