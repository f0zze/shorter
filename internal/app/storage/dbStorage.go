package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/f0zze/shorter/internal/app/entity"
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

func (d *PostgresStorage) Find(id string) (*ShortURL, bool) {

	query := `SELECT originalurl FROM urls WHERE shorturl = $1`

	result := &ShortURL{}
	err := d.db.QueryRowContext(context.Background(), query, id).Scan(&result.OriginalURL)

	if err != nil {
		return nil, false
	}

	return result, true
}

func (d *PostgresStorage) FindByUserID(id string) ([]entity.Shorter, error) {
	query := `SELECT * FROM urls WHERE userid = $1`

	rows, err := d.db.Query(query, id)

	if err != nil {
		return []entity.Shorter{}, err
	}

	defer rows.Close()

	var list []entity.Shorter
	for rows.Next() {
		item := entity.Shorter{}
		if err := rows.Scan(&item.ID, &item.ShortURL, &item.OriginalURL, &item.UserID); err != nil {
			return list, err
		}
		list = append(list, item)
	}

	if err = rows.Err(); err != nil {
		return list, err
	}

	return list, nil
}

func (d *PostgresStorage) Save(url []ShortURL) error {

	tx, err := d.db.Begin()

	if err != nil {
		return err
	}
	query := `INSERT INTO urls (id, shorturl, originalurl, userid) VALUES ($1, $2, $3, $4)`

	for _, u := range url {
		_, err := tx.ExecContext(context.Background(), query, u.UUID, u.ShortURL, u.OriginalURL, u.UserID)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
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

func (d *PostgresStorage) FindShortURLBy(originalURL string) (string, error) {

	var shortURL string
	query := `SELECT shorturl FROM urls WHERE originalurl = $1`

	err := d.db.QueryRowContext(context.Background(), query, originalURL).Scan(&shortURL)

	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (d *PostgresStorage) CreateTables() error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS urls (
			id VARCHAR(50) PRIMARY KEY,
			shortUrl VARCHAR(50) UNIQUE NOT NULL,
			originalUrl VARCHAR(50) UNIQUE NOT NULL,
			userId VARCHAR(50) NOT NULL
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
