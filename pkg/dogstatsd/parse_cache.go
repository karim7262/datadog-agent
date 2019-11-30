package dogstatsd

import (
	"strings"

	"github.com/arbll/ristretto"
)

type parserCache struct {
	tagsCache *ristretto.Cache
	nameCache *ristretto.Cache
}

func newParserCache() *parserCache {
	tagsCache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 2e6,
		MaxCost:     8388608, // 8MiB
		BufferItems: 64,      // number of keys per Get buffer.
	})
	nameCache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 2e6,
		MaxCost:     8388608, // 8MiB
		BufferItems: 64,      // number of keys per Get buffer.
	})

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
	c.tagsCache.Set(rawTags, tags, int64(len(rawTags)))
}

func (c *parserCache) GetName(rawName []byte) string {
	cachedName, found := c.nameCache.GetBytes(rawName)

	var name string
	if !found {
		name = string(rawName)
		c.putName(rawName, name)
	} else {
		name = cachedName.(string)
	}
	return name
}

func (c *parserCache) putName(rawName []byte, name string) {
	c.nameCache.Set(rawName, name, int64(len(name)))
}
