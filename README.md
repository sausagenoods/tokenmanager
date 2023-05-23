# tokenmanager
This is a token management library in Go.

## pgx (PostgreSQL) Store
### How to use

Create this table first:
```sql
CREATE TABLE IF NOT EXISTS tokens (
	token BYTEA NOT NULL,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL,
	owner TEXT NOT NULL
);
```

Example code using go-chi:
```go
package main

import (
        "context"
        "time"
        "fmt"
        "net/http"
        "log"

	"github.com/go-chi/chi/v5"
        "github.com/jackc/pgx/v5/pgxpool"
        tm "gitlab.com/sausagenoods/tokenmanager"
        "gitlab.com/sausagenoods/tokenmanager/store/pgxstore"
)

var tokenManager *tm.TokenManager

func init() {
        db, err := pgxpool.New(context.Background(), "postgresql://user:pass@localhost:5432/dbname")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
        }

        tokenManager = tm.New(pgxstore.New(db), nil)

	// Create an example token
	token, err := tm.GenerateToken()
        if err != nil {
                log.Fatal(err)
        }

	// Save the example token
        err := tokenManager.Commit(&tm.TokenInfo{
		Token: token,
		Owner: "Siren",
		Expiry: time.Now().Add(5 * time.Minute).UTC(),
		Data: []byte("Hello world")
	})
	if err != nil {
                log.Fatal(err)
        }

        log.Printf("Sample token is %s\n", token)
}

func main() {
        addr := ":3333"
        log.Printf("Starting server on %s\n", addr)
        http.ListenAndServe(addr, router())
}

func router() http.Handler {
        r := chi.NewRouter()

        // Protected routes
        r.Group(func(r chi.Router) {
                // Validate bearer tokens
                r.Use(tokenManager.TokenAuth)

                r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
                        data := tokenManager.FromContext(r.Context())
                        w.Write([]byte(fmt.Sprintf("Protected area. Hi %s\n.", data.Owner)))
                })
        })

        // Public routes
        r.Group(func(r chi.Router) {
                r.Get("/", func(w http.ResponseWriter, r *http.Request) {
                        w.Write([]byte("Public area."))
                })
        })

        return r
}
```
