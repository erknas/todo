package db

import (
	"fmt"
	"time"

	"github.com/zeze322/todo/lib"
	"github.com/zeze322/todo/repeattask"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type CreateTaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type CreateTaskResponse struct {
	ID string `json:"id"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

var mapping map[byte]bool = map[byte]bool{'d': true, 'y': true, 'w': true, 'm': true}

func NewTask(date, title, comment, repeat string) (Task, error) {
	if len(repeat) != 0 && !mapping[repeat[0]] {
		return Task{}, fmt.Errorf("unknown rule")
	}

	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().UTC().Location())

	if title == "" {
		return Task{}, fmt.Errorf("title should not be empty")
	}

	if date == "" {
		return Task{Date: now.Format(lib.Layout), Title: title, Comment: comment, Repeat: repeat}, nil
	}

	d, err := time.Parse(lib.Layout, date)
	if err != nil {
		return Task{}, fmt.Errorf("invalid date")
	}

	if repeat == "" {
		if !d.Before(now) {
			return Task{Date: date, Title: title, Comment: comment, Repeat: repeat}, nil
		} else {
			return Task{Date: now.Format(lib.Layout), Title: title, Comment: comment, Repeat: repeat}, nil
		}
	}

	switch repeat[0] {
	case 'd':
		next, err := repeattask.RepeatD(now, date, repeat)
		if err == nil {
			return Task{Date: next, Title: title, Comment: comment, Repeat: repeat}, nil
		} else {
			return Task{}, err
		}
	case 'y':
		next, err := repeattask.RepeatY(now, date, repeat)
		if err == nil {
			return Task{Date: next, Title: title, Comment: comment, Repeat: repeat}, nil
		} else {
			return Task{}, err
		}
	case 'w':
		next, err := repeattask.RepeatW(now, date, repeat)
		if err == nil {
			return Task{Date: next[0], Title: title, Comment: comment, Repeat: repeat}, nil
		} else {
			return Task{}, err
		}
	}
	return Task{}, nil
}
