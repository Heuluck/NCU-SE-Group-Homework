package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gophertodo/backend/internal/domain"
	"gophertodo/backend/internal/repository"
	"gophertodo/backend/internal/service"
)

type Server struct {
	tasks *service.TaskAppService
}

type taskInput struct {
	Content string `json:"content"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewServer(tasks *service.TaskAppService) *Server {
	return &Server{tasks: tasks}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/tasks", s.handleTasks)
	mux.HandleFunc("/tasks/", s.handleTaskByID)
	return withCORS(mux)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := s.tasks.ListTasks()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "failed to list tasks")
			return
		}
		writeJSON(w, http.StatusOK, tasks)
	case http.MethodPost:
		var input taskInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_input", "request body must be valid JSON")
			return
		}
		task, err := s.tasks.AddTask(input.Content)
		if err != nil {
			if errors.Is(err, domain.ErrEmptyContent) {
				writeError(w, http.StatusBadRequest, "invalid_input", "content is required")
				return
			}
			writeError(w, http.StatusInternalServerError, "internal_error", "failed to create task")
			return
		}
		writeJSON(w, http.StatusCreated, task)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	}
}

func (s *Server) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	id, action, ok := parseTaskPath(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}

	if action == "complete" {
		s.handleCompleteTask(w, r, id)
		return
	}
	if action != "" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, err := s.tasks.GetTask(id)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, task)
	case http.MethodDelete:
		if err := s.tasks.DeleteTask(id); err != nil {
			handleServiceError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, DELETE, OPTIONS")
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	}
}

func (s *Server) handleCompleteTask(w http.ResponseWriter, r *http.Request, id int) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	task, err := s.tasks.CompleteTask(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func parseTaskPath(path string) (id int, action string, ok bool) {
	rest := strings.TrimPrefix(path, "/tasks/")
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return 0, "", false
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id <= 0 {
		return 0, "", false
	}

	if len(parts) == 1 {
		return id, "", true
	}
	if len(parts) == 2 {
		return id, parts[1], true
	}
	return 0, "", false
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrTaskNotFound):
		writeError(w, http.StatusNotFound, "not_found", "task not found")
	case errors.Is(err, domain.ErrAlreadyDone):
		writeError(w, http.StatusConflict, "already_completed", "task is already completed")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "unexpected server error")
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, apiError{Code: code, Message: message})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
