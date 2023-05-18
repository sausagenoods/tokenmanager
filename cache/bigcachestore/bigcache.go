package bigcachestore

import (
	"context"

	"github.com/allegro/bigcache/v3"
)

type BigCacheStore struct {
	cache *BigCache
}

func (b *BigCacheStore) Delete(token []byte) (err error) {}
func (b *BigCacheStore) Find(token []byte) (b []byte, found bool, err error) {}
func (b *BigCacheStore) Commit(token []byte, b []byte, expiry time.Time) (err error) {}

// Does nothing because BigCache manages it's own evictions.
func (b *BigCacheStore) Cleanup() (err error) {
	return nil
}
