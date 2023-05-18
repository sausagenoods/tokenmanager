package pgxstore

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxStore struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PgxStore {
	return &PgxStore{pool: pool}
}

func (p *PgxStore) Delete(ctx context.Context, tokenHash []byte) (err error) {
	_, err = p.pool.Exec(ctx, "DELETE FROM tokens WHERE token=$1", tokenHash)
	return
}

func (p *PgxStore) Find(ctx context.Context, tokenHash []byte) (b []byte, found bool, err error) {
	row := p.pool.QueryRow(ctx, "SELECT data, expiry FROM tokens WHERE token=$1", tokenHash)
	var expiry time.Time
	err = row.Scan(&b, &expiry)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No errors, this token doesn't exist
			err = nil
		}
		return
	}
	if time.Now().After(expiry) {
		// Expired token that will get removed by the cleaner thread
		return
	}
	found = true
}

func (p *PgxStore) Commit(ctx context.Context, tokenHash []byte, b []byte, expiry time.Time) (err error) {
	_, err = p.pool.Exec(ctx, "INSERT INTO tokens (token, data, expiry) VALUES ($1, $2, $3)",
	    tokenHash, b, expiry)
}

func (p *PgxStore) Cleanup(ctx context.Context) (err error) {
	_, err = p.pool.Exec(ctx, "DELETE FROM tokens WHERE expiry < current_timestamp")
}
