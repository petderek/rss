package rss

import "errors"

func ToInternal(data []byte) (*Node, error) {
	// Try RSS first
	if rss, err := parseXML(data); err == nil {
		return FromRss(rss), nil
	}
	
	// Try Atom
	if atom, err := parseAtomXML(data); err == nil {
		return FromAtom(atom), nil
	}
	
	return nil, errors.New("unable to parse as RSS or Atom feed")
}
