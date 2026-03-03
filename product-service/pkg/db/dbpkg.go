package dbpkg

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	DB *pgx.Conn
}

func InitDB(dbUrl string) PostgresDB {
	db, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		panic("database does not work: " + err.Error())
	}
	return PostgresDB{DB: db}
}
