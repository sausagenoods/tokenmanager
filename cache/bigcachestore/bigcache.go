package bigcachestore

import (
	"context"

	"github.com/allegro/bigcache/v3"

	tm "gitlab.com/sausagenoods/tokenmanager"
)

type BigCacheStore struct {
	cache *BigCache
}

func (b *BigCacheStore) Delete(token string) error {}
func (b *BigCacheStore) Find(token string) (*tm.TokenInfo, error) {}
func (b *BigCacheStore) FindAllByOwner(owner string) ([]tm.TokenInfo, error) {}
func (b *BigCacheStore) Commit(info *tm.TokenInfo) error {}

// Does nothing because BigCache manages it's own evictions.
func (b *BigCacheStore) Cleanup() error {
	return nil
}
