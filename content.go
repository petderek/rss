package rss

import (
	"log"
)

type Content struct {
	Config SubscriptionConfig
	Cache  Cache
}

func (c Content) GetFeed(key string) ([]byte, error) {
	feed, err := c.Cache.GetSubscription(key)
	if err == nil { // todo: check ttl
		return feed.Data()
	}
	u, ok := c.Config[key]
	if !ok {
		log.Println("error: no config for key ", key)
	}
	data, err := fetch(u)
	if err != nil {
		return nil, err
	}
	if err = c.Cache.PutSubscription(key, data); err != nil {
		return nil, err
	}
	return data, nil
}
