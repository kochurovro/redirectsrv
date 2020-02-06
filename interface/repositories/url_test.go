package repositories

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/kochurovro/redirectsrv/domain"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Loader
)

func TestMain(m *testing.M) {
	db, err := sql.Open("mysql", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory("testdata/fixtures"),
	)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func TestSqlUrlRepo_Get(t *testing.T) {
	prepareTestDatabase()
	repo := NewSqlUrlRepo(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	defer func() {
		_, err := repo.DB.Exec("truncate table urls;")
		if err != nil {
			log.Println(err)
		}
	}()

	type compfunc func(got string) bool
	f := func(name string, compare compfunc, in string) {
		t.Run(name, func(t *testing.T) {
			out, err := repo.Get(in)
			if err != nil && err != domain.ErrRecordNotFound {
				t.Fatalf("expected nil error got %v", err)
				return
			}
			if !compare(out) {
				t.Error("")
				t.Log(out)
			}
		})
	}

	f("empty name should return", func(got string) bool { return got == "" }, "")
	f("one name should return firstURL", func(got string) bool { return got == "firstURL" }, "first")
}
