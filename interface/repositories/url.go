package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/kochurovro/redirectsrv/domain"
)

var _ domain.UrlGetter = &SqlUrlRepo{}

// SqlUrlRepo is simple mysql urls repository
type SqlUrlRepo struct {
	DB  *sql.DB
	ctx context.Context
}

func NewSqlUrlRepo(ctx context.Context, url string) *SqlUrlRepo {
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}
	return &SqlUrlRepo{DB: db, ctx: ctx}
}

// Get returns url and error by name
func (s *SqlUrlRepo) Get(name string) (url string, err error) {
	c := make(chan bool)
	defer close(c)

	go func() {
		for {
			err = s.DB.Ping()
			if err != nil {
				c <- true
			}
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		err = s.DB.QueryRow("SELECT url FROM urls WHERE name = ?", name).Scan(&url)
		c <- false
	}()
	connFail := <-c

	if err != nil {
		log.Println("repositories#get", err)
		if err == sql.ErrNoRows {
			return url, domain.ErrRecordNotFound
		}
		if connFail {
			return url, domain.ErrCloseConnection
		}
		return url, err
	}

	return url, nil
}
