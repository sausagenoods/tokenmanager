package tokenmanager

import (
	"log"
	"time"
)

type TokenManager struct {
	Store Store

	// A second token store for caching purposes, can be nil.
	Cache Store
}

type Store interface {
	Delete(token []byte) (err error)
	Find(token []byte) (b []byte, found bool, err error)
	Commit(token []byte, b []byte, expiry time.Time) (err error)
	Cleanup() (err error)
}

func New(store, cache Store) *TokenManager {
	return NewWithCleanupInterval(store, cache, 5*time.Minute)
}

func NewWithCleanupInterval(store Store, cache Store, interval time.Duration) *TokenManager {
	t := &TokenManager{Store: store, Cache: cache}
	if interval > 0 {
		go t.cleanup(interval)
	}
	return t
}

func (t *TokenManager) Find(token []byte) (b []byte, found bool, err error) {
	if t.Cache != nil {
		b, found, err = t.Cache.Find(token)
		if err != nil {
			return
		}
		if found {
			return
		}
	}
	return t.Store.Find(token)
}

func (t *TokenManager) Delete(token []byte) error {
	if t.Cache != nil {
		if err := t.Cache.Delete(token); err != nil {
			return err
		}
	}
	return t.Store.Delete(token)
}

func (t *TokenManager) Commit(token []byte, b []byte, expiry time.Time) error {
	if t.Cache != nil {
		if err := t.Cache.Commit(token, b, expiry); err != nil {
			return err
		}
	}
	return t.Store.Commit(token, b, expiry)
}

func (t *TokenManager) cleanup(interval time.Duration) {
	for {
		if t.Cache != nil {
			if err := t.Cache.Cleanup(); err != nil {
				log.Println(err)
			}
		}
		if err := t.Store.Cleanup(); err != nil {
			log.Println(err)
		}
		time.Sleep(interval)
	}
}
