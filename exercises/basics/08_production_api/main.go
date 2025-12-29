package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)


type task struct {
	Id int           `json:"id"`
	Name string      `json:"name"`
	Completed bool   `json:"completed"`
}


// We are using a map to store tasks in memory
// In go maps are not concurrent safe, so we use RWMutex to ensure concurrent access
type store struct {
	mu  sync.RWMutex
	tasks map[int]task
}


// We use (s *store) to pass the store instance to the function
// This is a pointer receiver, which means we are modifying the store instance
func (s *store) getTask(id int) (task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if id <= 0 {
		return task{}, fmt.Errorf("invalid task id: %d", id)
	}	
	if _, ok := s.tasks[id]; !ok {
		return task{}, fmt.Errorf("task %d not found", id)
	}
	return s.tasks[id], nil
}

func (s *store) getAllTasks() ([]task,error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	listTasks := make([]task,0,len(s.tasks))
	for _,task := range(s.tasks) {
		listTasks = append(listTasks, task)
	}
	return listTasks,nil
}

func (s *store) createTask(task task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task.Id <= 0 {
		return fmt.Errorf("invalid task id: %d", task.Id)
	}
	if _, ok := s.tasks[task.Id]; ok {
		return fmt.Errorf("task %d already exists", task.Id)
	}
	s.tasks[task.Id] = task
	return nil
}

func (s *store) updateTask(task task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task.Id <= 0 {
		return fmt.Errorf("invalid task id: %d", task.Id)
	}
	if _, ok := s.tasks[task.Id]; !ok {
		return fmt.Errorf("task %d not found", task.Id)
	}
	s.tasks[task.Id] = task
	return nil
}

func (s *store) deleteTask(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id <= 0 {
		return fmt.Errorf("invalid task id: %d", id)
	}
	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task %d not found", id)
	}
	delete(s.tasks, id)
	return nil
}

func WriteJson(w http.ResponseWriter, status int, v any)  {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err:= json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func main() {
	store := &store{tasks: make(map[int]task)}
	
	// CRUD server and endpoints
	mux:= http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Get all tasks
			allTasks,err := store.getAllTasks()
			if err != nil {
				WriteJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			WriteJson(w, http.StatusOK, allTasks)
		case http.MethodPost:
			// Create a new task
			var newTask task
			if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
				WriteJson(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
				return
			}
			if err := store.createTask(newTask); err != nil {
				WriteJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			WriteJson(w, http.StatusCreated, newTask)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Get a specific task
			var id int
			if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
				WriteJson(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
				return
			}
			task,err := store.getTask(id)
			if err != nil {
				WriteJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			WriteJson(w, http.StatusOK, task)
		case http.MethodPut:
			// Update a task
			var updatedTask task
			if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
				WriteJson(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
				return
			}
			if err := store.updateTask(updatedTask); err != nil {
				WriteJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			WriteJson(w, http.StatusOK, updatedTask)
		case http.MethodDelete:
			// Delete a task
			var id int
			if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
				WriteJson(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
				return
			}
			if err := store.deleteTask(id); err != nil {
				WriteJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			WriteJson(w, http.StatusOK, map[string]string{"message": "task deleted successfully"})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})


	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	go func() {
		fmt.Println("Server started on :8080")
		if err:= server.ListenAndServe(); err != nil {
			fmt.Println("Server error:", err)
		}
	} ()
}
