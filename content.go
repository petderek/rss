package rss

import (
	"fmt"
	"log"
)

type Content struct {
	Config SubscriptionConfig
	Cache  Cache
}

// GetFeed returns raw RSS data for a feed
func (c Content) GetFeed(key string) ([]byte, error) {
	feed, err := c.Cache.GetSubscription(key)
	if err == nil { // todo: check ttl
		return feed.Data()
	}
	u, ok := c.Config[key]
	if !ok {
		log.Println("error: no config for key ", key)
		return nil, fmt.Errorf("no config for feed: %s", key)
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

// GetNode returns parsed RSS feed as filesystem node
func (c Content) GetNode(feedName string) (*Node, error) {
	data, err := c.GetFeed(feedName)
	if err != nil {
		return nil, err
	}
	node, err := ToInternal(data)
	if err != nil {
		return nil, err
	}
	node.Name = feedName
	return node, nil
}

// ListFeeds returns all configured feed names
func (c Content) ListFeeds() []string {
	feeds := make([]string, 0, len(c.Config))
	for name := range c.Config {
		feeds = append(feeds, name)
	}
	return feeds
}

// RefreshFeed forces a refresh of the feed from the network
func (c Content) RefreshFeed(feedName string) error {
	u, ok := c.Config[feedName]
	if !ok {
		return fmt.Errorf("no config for feed: %s", feedName)
	}
	data, err := fetch(u)
	if err != nil {
		return err
	}
	return c.Cache.PutSubscription(feedName, data)
}
