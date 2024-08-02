package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zeze322/todo/lib"
)

const limit = 25

type Storage interface {
	CreateTask(context.Context, Task) (string, error)
	GetTasks(context.Context, string) ([]Task, error)
	GetTask(context.Context, string) (Task, error)
	UpdateTask(context.Context, string, Task) error
	DeleteTask(context.Context, string) error
}

type SqliteStorage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*SqliteStorage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), storagePath)

	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	if install {
		stmtTable, err := db.Prepare(`
			CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date INTEGER,
			title TEXT,
			comment TEXT,
			repeat TEXT CHECK(LENGTH(repeat) <= 128)
			);
		`)

		if err != nil {
			return nil, err
		}

		defer stmtTable.Close()

		_, err = stmtTable.Exec()
		if err != nil {
			return nil, err
		}

		stmtIdxDate, err := db.Prepare(`
			CREATE INDEX IF NOT EXISTS idx_date on scheduler (date);
		`)
		if err != nil {
			return nil, err
		}

		defer stmtIdxDate.Close()

		_, err = stmtIdxDate.Exec()
		if err != nil {
			return nil, err
		}

	}

	return &SqliteStorage{
		db: db,
	}, nil
}

func (s *SqliteStorage) Close() error {
	return s.db.Close()
}

func (s *SqliteStorage) CreateTask(ctx context.Context, task Task) (string, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)`

	res, err := s.db.ExecContext(ctx, query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(id)), nil
}

func (s *SqliteStorage) GetTasks(ctx context.Context, keyWord string) ([]Task, error) {
	if keyWord == "" {
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT $1`

		rows, err := s.db.Query(query, limit)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		var tasks []Task

		for rows.Next() {
			task := Task{}
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
		}

		if len(tasks) == 0 {
			return []Task{}, nil
		}

		return tasks, nil
	}

	if keyWord != "" && lib.IsDate(keyWord) {
		date, err := lib.ParseTime(keyWord)
		if err != nil {
			return nil, err
		}

		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date=$1`

		rows, err := s.db.Query(query, date)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		var tasks []Task

		for rows.Next() {
			task := Task{}
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
		}

		if len(tasks) == 0 {
			return []Task{}, nil
		}

		return tasks, nil
	}

	if keyWord != "" && !lib.IsDate(keyWord) {
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE LOWER(title) LIKE $1 OR LOWER(comment) LIKE $1`
		keyWord = "%" + keyWord + "%"

		rows, err := s.db.Query(query, keyWord)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		var tasks []Task

		for rows.Next() {
			task := Task{}
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
		}

		if len(tasks) == 0 {
			return []Task{}, nil
		}

		return tasks, nil
	}

	return nil, nil
}

func (s *SqliteStorage) GetTask(ctx context.Context, id string) (Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id=$1`
	row := s.db.QueryRowContext(ctx, query, id)

	var task Task

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, fmt.Errorf("task not found id: %s", id)
	} else if err != nil {
		return Task{}, fmt.Errorf("failed to get task")
	}

	return task, nil
}

func (s *SqliteStorage) UpdateTask(ctx context.Context, id string, task Task) error {
	query := `UPDATE scheduler SET date=$1, title=$2, comment=$3, repeat=$4 WHERE id=$5`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to update task")
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, task.Date, task.Title, task.Comment, task.Repeat, id)
	if err != nil {
		return fmt.Errorf("failed to update task")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("task not found id: %s", id)
	}

	return nil
}

func (s *SqliteStorage) DeleteTask(ctx context.Context, id string) error {
	query := `DELETE FROM scheduler WHERE id=$1`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete task")
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete task")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("task not found id: %s", id)
	}

	return nil
}
