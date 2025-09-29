package rss

import (
	"bufio"
	"bytes"
	"net/url"
	"strings"
)

type NormalConfig map[string]string

type SubscriptionConfig map[string]*url.URL

func FromFile(filename string) (NormalConfig, error) {
	config := map[string]string{}
	data, err := readFile(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), "=")
		if len(tokens) != 2 {
			continue
		}
		left := strings.TrimSpace(tokens[0])
		right := strings.TrimSpace(tokens[1])
		config[left] = right
	}
	return config, nil
}

func ToSubscription(nc NormalConfig) SubscriptionConfig {
	urls := map[string]*url.URL{}
	for k, v := range nc {
		parsed, err := url.Parse(v)
		if err != nil {
			continue
		}
		urls[k] = parsed
	}
	return urls
}

func (sc SubscriptionConfig) String() string {
	sb := strings.Builder{}
	for k, v := range sc {
		if sb.Len() > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(k)
		sb.WriteRune('=')
		sb.WriteString(v.String())
	}
	return sb.String()
}
