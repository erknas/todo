package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/zeze322/todo/db"
	"github.com/zeze322/todo/lib"
	"github.com/zeze322/todo/repeattask"
)

func (s *Server) handleNextDate(w http.ResponseWriter, r *http.Request) error {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	parseNow, err := time.Parse(lib.Layout, now)
	if err != nil {
		return err
	}

	if len(repeat) == 0 {
		return lib.WriteJSON(w, http.StatusBadRequest, lib.ApiErr{Error: error.Error(errors.New("empty rule"))})
	}

	switch repeat[0] {
	case 'd':
		next, err := repeattask.RepeatD(parseNow, date, repeat)
		if err != nil {
			return err
		}
		return lib.WriteJSON(w, http.StatusOK, next)
	case 'y':
		next, err := repeattask.RepeatY(parseNow, date, repeat)
		if err != nil {
			return err
		}
		lib.WriteJSON(w, http.StatusOK, next)
	case 'w':
		next, err := repeattask.RepeatW(parseNow, date, repeat)
		if err != nil {
			return err
		}
		return lib.WriteJSON(w, http.StatusOK, next)
	case 'm':
		next, err := repeattask.RepeatM(parseNow, date, repeat)
		if err != nil {
			return err
		}
		return lib.WriteJSON(w, http.StatusOK, next)
	default:
		return lib.WriteJSON(w, http.StatusBadRequest, lib.ApiErr{Error: "empty rule"})
	}

	return nil
}

func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetTasks(w, r)
	case "DELETE":
		return s.handleDeleteTask(w, r)
	case "POST":
		return s.handleCreateTask(w, r)
	case "PUT":
		return s.handleUpdateTask(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) error {
	var req db.CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	task, err := db.NewTask(req.Date, req.Title, req.Comment, req.Repeat)
	if err != nil {
		return err
	}

	id, err := s.store.CreateTask(r.Context(), task)
	if err != nil {
		return err
	}

	return lib.WriteJSON(w, http.StatusOK, db.CreateTaskResponse{
		ID: id,
	})
}

func (s *Server) handleGetTasks(w http.ResponseWriter, r *http.Request) error {
	search := r.FormValue("search")

	tasks, err := s.store.GetTasks(r.Context(), search)
	if err != nil {
		return fmt.Errorf("failed to get tasks")
	}

	return lib.WriteJSON(w, http.StatusOK, db.TasksResponse{Tasks: tasks})
}

func (s *Server) handleDeleteTask(w http.ResponseWriter, r *http.Request) error {
	id := r.FormValue("id")
	if id == "" {
		return fmt.Errorf("id not specified")
	}

	if err := s.store.DeleteTask(r.Context(), id); err != nil {
		return err
	}

	return lib.WriteJSON(w, http.StatusOK, lib.EmptyJSON{})
}

func (s *Server) handleGetTaskByID(w http.ResponseWriter, r *http.Request) error {
	id := r.FormValue("id")
	if id == "" {
		return fmt.Errorf("id not specified")
	}

	task, err := s.store.GetTask(r.Context(), id)
	if err != nil {
		return err
	}

	return lib.WriteJSON(w, http.StatusOK, task)
}

func (s *Server) handleUpdateTask(w http.ResponseWriter, r *http.Request) error {
	var req db.Task

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	updateTask, err := db.NewTask(req.Date, req.Title, req.Comment, req.Repeat)
	if err != nil {
		return err
	}

	if err := s.store.UpdateTask(r.Context(), req.ID, updateTask); err != nil {
		return err
	}

	return lib.WriteJSON(w, http.StatusOK, lib.EmptyJSON{})
}

func (s *Server) handleTaskDone(w http.ResponseWriter, r *http.Request) error {
	id := r.FormValue("id")
	if id == "" {
		return fmt.Errorf("id not specified")
	}

	task, err := s.store.GetTask(r.Context(), id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		if err := s.store.DeleteTask(r.Context(), id); err != nil {
			return err
		}
	}

	if task.Repeat != "" {
		next, err := repeattask.UpdateDate(task.Date, task.Repeat)
		if err != nil {
			return err
		}

		task.Date = next

		if err := s.store.UpdateTask(r.Context(), id, task); err != nil {
			return err
		}
	}

	return lib.WriteJSON(w, http.StatusOK, lib.EmptyJSON{})
}
