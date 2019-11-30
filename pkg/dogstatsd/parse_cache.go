package dogstatsd

import (
	"strings"

	"github.com/DataDog/datadog-agent/pkg/util/cache"
)

type parserCache struct {
	tagsCache *cache.LRUCache
	nameCache *cache.LRUCache
}

func newParserCache() *parserCache {
	tagsCache := cache.NewLRUCache(1024)
	nameCache := cache.NewLRUCache(1024)

	return &parserCache{
		tagsCache: tagsCache,
		nameCache: nameCache,
	}
}

func (c *parserCache) GetTags(rawTags []byte) []string {
	var tags []string
	cachedTags, found := c.tagsCache.GetBytes(rawTags)
	if !found {
		tags = strings.Split(string(rawTags), commaSeparatorString)
		c.putTags(rawTags, tags)
	} else {
		tags = cachedTags.([]string)
	}
	return tags
}

func (c *parserCache) putTags(rawTags []byte, tags []string) {
	c.tagsCache.Put(string(rawTags), tags)
}

func (c *parserCache) GetName(rawName []byte) string {
	cachedName, found := c.nameCache.GetBytes(rawName)

	var name string
	if !found {
		name = string(rawName)
		c.putName(name)
	} else {
		name = cachedName.(string)
	}
	return name
}

func (c *parserCache) putName(name string) {
	c.nameCache.Put(name, name)
}
