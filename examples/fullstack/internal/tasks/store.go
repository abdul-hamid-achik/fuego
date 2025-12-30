// Package tasks provides a shared in-memory task store for the fullstack example.
package tasks

import (
	"strconv"
	"sync"
)

// Task represents a simple task item.
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Store manages tasks in memory with thread-safe operations.
type Store struct {
	tasks  []Task
	mu     sync.RWMutex
	nextID int
}

// Default is the shared task store instance used by all handlers.
var Default = &Store{nextID: 1}

// List returns all tasks.
func (s *Store) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external mutation
	result := make([]Task, len(s.tasks))
	copy(result, s.tasks)
	return result
}

// Add creates a new task with the given title.
func (s *Store) Add(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:        s.nextID,
		Title:     title,
		Completed: false,
	}
	s.tasks = append(s.tasks, task)
	s.nextID++
	return task
}

// Delete removes a task by ID.
func (s *Store) Delete(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return
		}
	}
}

// Toggle switches the completed status of a task by ID.
func (s *Store) Toggle(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].Completed = !task.Completed
			return
		}
	}
}

// RenderHTML returns the task list as HTML for HTMX responses.
func (s *Store) RenderHTML() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.tasks) == 0 {
		return `<p class="text-gray-500 italic">No tasks yet. Add one above!</p>`
	}

	html := `<ul class="space-y-2">`
	for _, task := range s.tasks {
		checkedClass := ""
		checkedAttr := ""
		if task.Completed {
			checkedClass = " line-through text-gray-400"
			checkedAttr = ` checked`
		}
		html += `<li class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
			<input 
				type="checkbox" 
				class="h-4 w-4 text-orange-600 focus:ring-orange-500 border-gray-300 rounded"
				hx-post="/api/tasks/toggle?id=` + strconv.Itoa(task.ID) + `"
				hx-target="#task-list"
				hx-swap="innerHTML"` + checkedAttr + `
			/>
			<span class="flex-1` + checkedClass + `">` + task.Title + `</span>
			<button 
				class="text-red-500 hover:text-red-700"
				hx-delete="/api/tasks?id=` + strconv.Itoa(task.ID) + `"
				hx-target="#task-list"
				hx-swap="innerHTML"
			>
				Delete
			</button>
		</li>`
	}
	html += `</ul>`

	return html
}
