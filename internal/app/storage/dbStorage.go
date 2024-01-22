package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	db *sql.DB
}

func connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

//type Storage interface {
//	Find(uuid string) (*ShortURL, bool)
//	Save(url *ShortURL) error
//	Size() int
//	Ping() bool
//}

func (d *PostgresStorage) Find(_ string) (*ShortURL, bool) {
	/*For now do noting*/

	return nil, false
}

func (d *PostgresStorage) Save(_ *ShortURL) error {
	/*For now do noting*/
	return nil
}

func (d *PostgresStorage) Size() int {
	return 0
}

func (d *PostgresStorage) Ping() bool {
	err := d.db.Ping()

	return err != nil
}

func NewPostgresStorage(dsn string) (Storage, error) {
	db, err := connect(dsn)

	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db,
	}, nil
}
