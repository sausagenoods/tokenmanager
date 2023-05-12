package bigcachestore

import (
	"context"

	"github.com/allegro/bigcache/v3"
)

type BigCacheStore struct {
	cache *BigCache
}

func (b *BigCacheStore) Delete(token string) (err error) {}
func (b *BigCacheStore) Find(token string) (b []byte, found bool, err error) {}
func (b *BigCacheStore) Commit(token string, b []byte, expiry time.Time) (err error) {}
