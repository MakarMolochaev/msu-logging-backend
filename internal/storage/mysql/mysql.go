package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	const op = "storage.mysql.New"

	connStr := os.Getenv("MYSQL_CONN_STR")
	if connStr == "" {
		return nil, fmt.Errorf("%s: missing MySQL connection string in MYSQL_CONN_STR environment variable", op)
	}

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open db connection: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: failed to ping db: %w", op, err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return &Storage{db: db}, nil
}

func (s *Storage) SaveAudioFile(ctx context.Context, link string) (int64, error) {
	const op = "storage.mysql.SaveAudioFile"

	stmt, err := s.db.Prepare("INSERT INTO logging.audio_file (link, date_created) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

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
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

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
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

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
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

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

func (s *Storage) CreateNewTaskStatus(ctx context.Context) (int64, error) {
	const op = "storage.mysql.CreateNewTaskStatus"
	var task_status string = "none"

	stmt, err := s.db.Prepare("INSERT INTO logging.tasks (task_status) VALUES (?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, task_status)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) UpdateTaskStatusByID(ctx context.Context, id int64, newStatus string) error {
	const op = "storage.mysql.UpdateTaskStatusByID"

	stmt, err := s.db.Prepare("UPDATE logging.tasks SET task_status = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, newStatus, id)
	if err != nil {
		return fmt.Errorf("%s: execute query: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: no task found with id %d", op, id)
	}

	return nil
}

func (s *Storage) GetTaskStatusByID(ctx context.Context, id int64) (string, error) {
	const op = "storage.mysql.GetTaskStatusByID"

	stmt, err := s.db.Prepare("SELECT task_status FROM logging.tasks WHERE id = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var taskStatus string

	err = stmt.QueryRowContext(ctx, id).Scan(&taskStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: task with id %d not found", op, id)
		}
		return "", fmt.Errorf("%s: execute query: %w", op, err)
	}

	return taskStatus, nil
}
