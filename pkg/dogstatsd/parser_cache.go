package dogstatsd

import (
	"github.com/dgraph-io/ristretto"
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

func (c *parserCache) GetTags(rawTags []byte) ([]string, bool) {
	tags, found := c.tagsCache.Get(rawTags)
	return tags.([]string), found
}

func (c *parserCache) PutTags(rawTags []byte, tags []string, tagsSize int64) {
	c.tagsCache.Set(rawTags, tags, tagsSize)
}

func (c *parserCache) GetName(rawName []byte) (string, bool) {
	name, found := c.tagsCache.Get(rawName)
	return name.(string), found
}

func (c *parserCache) PutName(rawName []byte, name string) {
	c.tagsCache.Set(rawName, name, int64(len(name)))
}
