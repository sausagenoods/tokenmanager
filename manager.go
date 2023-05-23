package tokenmanager

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"log"
	"time"
)

type TokenManager struct {
	Store Store

	// A second token store for caching purposes, can be nil.
	Cache Store

	// Key used to store token data in router context
	contextKey string
}

// Represents a token and its metadata
type TokenInfo struct {
	Token  string
	Data   []byte
	Owner  string
	Expiry time.Time
}

// Token store methods
type Store interface {
	Delete(token string) (err error)
	Find(token string) (info *TokenInfo, err error)
	FindAllByOwner(owner string) (info []TokenInfo, err error)
	Commit(info *TokenInfo) (err error)
	Cleanup() (err error)
}

var (
	ErrExpired  = errors.New("Token is expired")
	ErrNotFound = errors.New("Token not found")
)

func New(store, cache Store) *TokenManager {
	return NewWithCleanupInterval(store, cache, 5*time.Minute)
}

func NewWithCleanupInterval(store Store, cache Store, interval time.Duration) *TokenManager {
	t := &TokenManager{Store: store, Cache: cache, contextKey: generateContextKey()}
	if interval > 0 {
		go t.cleanup(interval)
	}
	return t
}

func (t *TokenManager) Find(token string) (*TokenInfo, error) {
	if t.Cache != nil {
		info, err := t.Cache.Find(token)
		if err != nil {
			return info, err
		}
	}
	return t.Store.Find(token)
}

func (t *TokenManager) FindAllByOwner(owner string) ([]TokenInfo, error) {
	if t.Cache != nil {
		info, err := t.Cache.FindAllByOwner(owner)
		if err != nil {
			return info, err
		}
	}
	return t.Store.FindAllByOwner(owner)
}

func (t *TokenManager) Delete(token string) error {
	if t.Cache != nil {
		if err := t.Cache.Delete(token); err != nil {
			return err
		}
	}
	return t.Store.Delete(token)
}

func (t *TokenManager) Commit(info *TokenInfo) error {
	if t.Cache != nil {
		if err := t.Cache.Commit(info); err != nil {
			return err
		}
	}
	return t.Store.Commit(info)
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

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).
		EncodeToString(b), nil
}
