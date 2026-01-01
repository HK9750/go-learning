package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	Name        string `json:"name"`
	Id          int    `json:"id"`
	IsCompleted bool   `json:"isCompleted"`
}

type store struct {
	mu     sync.RWMutex
	tasks  map[int]Task
	nextId int
}

func GetStore() *store {
	return &store{
		tasks:  make(map[int]Task),
		nextId: 0,
	}
}

func GetServer(mux *http.ServeMux, port string) *http.Server {
	return &http.Server{
		Addr:    port,
		Handler: mux,
	}
}

func (s *store) getAllTasks() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *store) createTask(name string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextId++
	t := Task{
		Id:          s.nextId,
		Name:        name,
		IsCompleted: false,
	}
	s.tasks[t.Id] = t
	return t
}

func (s *store) getTaskById(id int) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if id <= 0 {
		return Task{}, fmt.Errorf("id should be greater than zero")
	}
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task not found")
	}
	return task, nil
}

func (s *store) updateTask(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[task.Id]; !ok {
		return fmt.Errorf("task not found")
	}
	s.tasks[task.Id] = task
	return nil
}

func (s *store) deleteTask(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task not found")
	}
	delete(s.tasks, id)
	return nil
}


func (s *store) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks := s.getAllTasks()
		WriteJson(w, http.StatusOK, tasks)
	case http.MethodPost:
		var input struct {
			Name string `json:"name"`
		}
		if err := ReadJson(r, &input); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid JSON payload")
			return
		}
		if strings.TrimSpace(input.Name) == "" {
			WriteError(w, http.StatusBadRequest, "name cannot be empty")
			return
		}
		newTask := s.createTask(input.Name)
		WriteJson(w, http.StatusCreated, newTask)
	default:
		w.Header().Set("Allow", "GET, POST")
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *store) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromPath(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, err := s.getTaskById(id)
		if err != nil {
			WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		WriteJson(w, http.StatusOK, task)

	case http.MethodPut:
		var t Task
		if err := ReadJson(r, &t); err != nil {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		t.Id = id
		if err := s.updateTask(t); err != nil {
			WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		WriteJson(w, http.StatusOK, t)

	case http.MethodDelete:
		if err := s.deleteTask(id); err != nil {
			WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "error", err)
		http.Error(w, `{"error" : "internal server error"}`, http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJson(w, status, map[string]string{"error": message})
}

func ReadJson(r *http.Request, t any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(t)
}

func getIdFromPath(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}


func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
			"remote_addr", r.RemoteAddr,
		)
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":5000"
	} else if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	store := GetStore()
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", store.handleTasks)
	mux.HandleFunc("/tasks/{id}", store.handleTaskByID)

	server := GetServer(mux, port) 
	handler := LoggerMiddleware(mux)
	// handler := chain(mux,)
	server.Handler = handler

	go func() {
		slog.Info("server starting", "addr", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	slog.Info("server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	} else {
		slog.Info("server stopped gracefully")
	}
}
