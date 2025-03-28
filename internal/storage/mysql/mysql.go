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

func (s *Storage) SaveAudioFile(ctx context.Context, link string) (int64, error) {
	const op = "storage.mysql.SaveAudioFile"

	stmt, err := s.db.Prepare("INSERT INTO logging.audio_file (link, date_created) VALUES (?, ?)")

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

func (s *Storage) SaveTextFile(ctx context.Context, text_full, text_short string) (int64, error) {
	const op = "storage.mysql.SaveTextFile"

	stmt, err := s.db.Prepare("INSERT INTO logging.text_file (text_full, text_short, date_created) VALUES (?, ?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, text_full, text_short, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) SaveValuation(ctx context.Context, usability, processing_speed, processing_quality int, reuse_service bool, comment string) (int64, error) {
	const op = "storage.mysql.SaveValuation"

	stmt, err := s.db.Prepare("INSERT INTO logging.valuation (usability, processing_speed, processing_quality, reuse_service, comment) VALUES (?, ?, ?, ?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, usability, processing_speed, processing_quality, reuse_service, comment)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) SaveMetrics(ctx context.Context, image_count, guest_count int, av_audio_time, av_process_time float32, satisfy_user_count int) (int64, error) {
	const op = "storage.mysql.SaveMetrics"

	stmt, err := s.db.Prepare("INSERT INTO logging.metrics (image_count, guest_count, av_audio_time, av_process_time, satisfy_user_count) VALUES (?, ?, ?, ?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, image_count, guest_count, av_audio_time, av_process_time, satisfy_user_count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
