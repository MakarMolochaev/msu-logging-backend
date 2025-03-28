package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	const op = "storage.mysql.New"

	db, err := sql.Open("mysql", os.Getenv("MYSQL_CONN_STR"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveAudioFileLink(ctx context.Context, link string) (int64, error) {
	const op = "storage.mysql.SaveAudioFileLink"

	stmt, err := s.db.Prepare("INSERT INTO logging.audio_file (link, date_au) VALUES (?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, link, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
