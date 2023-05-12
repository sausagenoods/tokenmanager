package tokenmanager

type TokenManager struct {
        Store Store
	Cache Cache
}

type Store interface {
        Delete(token string) (err error)
        Find(token string) (b []byte, found bool, err error)
        Commit(token string, b []byte, expiry time.Time) (err error)
}

type Cache interface {
	Store
}
