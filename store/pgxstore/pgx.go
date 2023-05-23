package pgxstore

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	tm "gitlab.com/sausagenoods/tokenmanager"
)

type PgxStore struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PgxStore {
	return &PgxStore{pool: pool}
}

func (p *PgxStore) Delete(token string) error {
	_, err := p.pool.Exec(context.Background(), "DELETE FROM tokens WHERE token=$1", token)
	return err
}

func (p *PgxStore) Find(token string) (*tm.TokenInfo, error) {
	row := p.pool.QueryRow(context.Background(),
		"SELECT data,owner,expiry FROM tokens WHERE token=$1", token)

	info := &tm.TokenInfo{Token: token}
	if err := row.Scan(info.Data, info.Owner, info.Expiry); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No errors, this token doesn't exist
			return nil, tm.ErrNotFound
		}
		return nil, err
	}
	if time.Now().After(info.Expiry) {
		// Expired token that will get removed by the cleaner thread
		return nil, tm.ErrExpired
	}
	return info, nil
}

func (p *PgxStore) FindAllByOwner(owner string) ([]tm.TokenInfo, error) {
	rows, err := p.pool.Query(context.Background(),
		"SELECT token,data,expiry from tokens WHERE owner=$1", owner)
	if err != nil {
		return nil, err
	}

	var tokens []tm.TokenInfo
	for rows.Next() {
		t := tm.TokenInfo{Owner: owner}
		if err := rows.Scan(&t.Token, &t.Data, &t.Expiry); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func (p *PgxStore) Commit(info *tm.TokenInfo) error {
	_, err := p.pool.Exec(context.Background(),
		"INSERT INTO tokens(token,data,owner,expiry)VALUES($1,$2,$4,$3)",
		info.Token, info.Data, info.Owner, info.Expiry)
	return err
}

func (p *PgxStore) Cleanup() error {
	_, err := p.pool.Exec(context.Background(),
		"DELETE FROM tokens WHERE expiry<current_timestamp")
	return err
}
