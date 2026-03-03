package dbpkg

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func InitDB(dbUrl string) PostgresDB {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic("database does not work: " + err.Error())
	}
	return PostgresDB{DB: db}
}
