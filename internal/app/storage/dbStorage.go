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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (d *PostgresStorage) Find(_ string) (*ShortURL, bool) {
	/*For now do noting*/

	return nil, false
}

func (d *PostgresStorage) Save(url []ShortURL) error {

	for _, u := range url {
		_, err := d.db.ExecContext(context.Background(), `
		INSERT INTO urls (id, shorturl, originalurl)
		VALUES ($1, $2, $3)
    `, u.UUID, u.ShortURL, u.OriginalURL)

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *PostgresStorage) Size() int {
	return 0
}

func (d *PostgresStorage) Ping() bool {
	err := d.db.Ping()

	return err == nil
}

func (d *PostgresStorage) Close() error {
	return d.db.Close()
}

func (d *PostgresStorage) CreateTables() error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS urls (
			id VARCHAR(50) PRIMARY KEY,
			shortUrl VARCHAR(50) UNIQUE NOT NULL,
			originalUrl VARCHAR(50) NOT NULL
		)
	`

	_, err := d.db.Exec(createTableQuery)

	return err
}

func NewPostgresStorage(dsn string) (Storage, error) {
	db, err := connect(dsn)

	if err != nil {
		return nil, err
	}

	storage := &PostgresStorage{
		db,
	}

	err = storage.CreateTables()

	if err != nil {
		return nil, err
	}

	return storage, nil
}
