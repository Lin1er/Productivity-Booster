package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "Low"
	case PriorityMedium:
		return "Medium"
	case PriorityHigh:
		return "High"
	}
	return "Unknown"
}

type Todo struct {
	ID          int
	Title       string
	Description string
	Completed   bool
	Priority    Priority
	CreatedAt   time.Time
	DueTime     *time.Time
}

type TodoList struct {
	db       *sql.DB
	Todos    []*Todo
	Selected int
	NextID   int
}

// Load - Load semua todos dari database ke memory
func (tl *TodoList) Load() error {
	query := "SELECT id, title, description, completed, priority, created_at, due_date FROM todos ORDER BY id"
	rows, err := tl.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	tl.Todos = []*Todo{} // Clear existing

	for rows.Next() {
		var id int
		var title, description string
		var completed bool
		var priority int
		var createdAt time.Time
		var dueDate sql.NullTime

		if err := rows.Scan(&id, &title, &description, &completed, &priority, &createdAt, &dueDate); err != nil {
			return fmt.Errorf("failed to scan todo: %w", err)
		}

		todo := &Todo{
			ID:          id,
			Title:       title,
			Description: description,
			Completed:   completed,
			Priority:    Priority(priority),
			CreatedAt:   createdAt,
		}

		if dueDate.Valid {
			todo.DueTime = &dueDate.Time
		}

		tl.Todos = append(tl.Todos, todo)

		// Update NextID
		if id >= tl.NextID {
			tl.NextID = id + 1
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating todos: %w", err)
	}

	return nil
}

/*
* testing purpose
* NewTodoListSample creates an empty TodoList for
 */

func NewTodoListSample() *TodoList {
	return &TodoList{
		Todos: []*Todo{
			{
				ID: 1, Title: "Buy groceries",
				Description: "Milk, Bread, Eggs",
				Completed:   false,
				Priority:    PriorityMedium,
				CreatedAt:   time.Now(),
			},
			{
				ID:          2,
				Title:       "Finish project",
				Description: "Complete the Go project by Friday",
				Completed:   false,
				Priority:    PriorityHigh,
				CreatedAt:   time.Now(),
			},
			{
				ID:          3,
				Title:       "Call Mom",
				Description: "Check in with Mom this weekend",
				Completed:   true,
				Priority:    PriorityLow,
				CreatedAt:   time.Now(),
			},
			{
				ID:          4,
				Title:       "Workout",
				Description: "Go for a 30-minute run",
				Completed:   false,
				Priority:    PriorityMedium,
				CreatedAt:   time.Now(),
			},
			{
				ID:          5,
				Title:       "Read a book",
				Description: "Finish reading 'The Go Programming Language'",
				Completed:   true,
				Priority:    PriorityLow,
				CreatedAt:   time.Now(),
			},
			{
				ID:          6,
				Title:       "Plan vacation",
				Description: "Research destinations for summer vacation",
				Completed:   false,
				Priority:    PriorityHigh,
				CreatedAt:   time.Now(),
			},
		},
		Selected: 0,
		NextID:   7,
	}
}

func NewTodoList(db_ *sql.DB) *TodoList {
	tl := &TodoList{
		db:       db_,
		Todos:    []*Todo{},
		Selected: 0,
		NextID:   1,
	}
	// Auto-load dari database saat inisialisasi
	if err := tl.Load(); err != nil {
		// Log error tapi tetap return instance kosong
		fmt.Printf("Warning: failed to load todos: %v\n", err)
	}
	return tl
}

func (tl *TodoList) SelectNext() {
	if tl.Selected < len(tl.Todos)-1 {
		tl.Selected++
	}
}

func (tl *TodoList) SelectPrevious() {
	if tl.Selected > 0 {
		tl.Selected--
	}
}

func (tl *TodoList) GetSelected() *Todo {
	if len(tl.Todos) == 0 || tl.Selected >= len(tl.Todos) {
		return nil
	}
	return tl.Todos[tl.Selected]
}

func (tl *TodoList) Count() int {
	return len(tl.Todos)
}

// Add - Tambah todo ke database DAN memory sekaligus
func (tl *TodoList) Add(title, description string, priority Priority, dueTime *time.Time) error {
	query := `INSERT INTO todos (title, description, completed, priority, due_date, created_at)
	          VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := tl.db.Exec(query, title, description, false, int(priority), dueTime, now)
	if err != nil {
		return fmt.Errorf("failed to add todo to database: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Tambah ke memory
	todo := &Todo{
		ID:          int(id),
		Title:       title,
		Description: description,
		Completed:   false,
		Priority:    priority,
		CreatedAt:   now,
		DueTime:     dueTime,
	}
	tl.Todos = append(tl.Todos, todo)
	tl.NextID = int(id) + 1

	return nil
}

// Update - Update todo di database DAN memory sekaligus
func (tl *TodoList) Update(id int, title, description string, priority Priority, dueTime *time.Time) error {
	query := `UPDATE todos SET title=?, description=?, priority=?, due_date=?, updated_at=CURRENT_TIMESTAMP
	          WHERE id=?`

	_, err := tl.db.Exec(query, title, description, int(priority), dueTime, id)
	if err != nil {
		return fmt.Errorf("failed to update todo in database: %w", err)
	}

	// Update di memory
	for _, todo := range tl.Todos {
		if todo.ID == id {
			todo.Title = title
			todo.Description = description
			todo.Priority = priority
			todo.DueTime = dueTime
			break
		}
	}

	return nil
}

// Remove - Hapus todo dari database DAN memory sekaligus
func (tl *TodoList) Remove(id int) error {
	query := `DELETE FROM todos WHERE id=?`

	_, err := tl.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo from database: %w", err)
	}

	// Hapus dari memory
	for i, todo := range tl.Todos {
		if todo.ID == id {
			tl.Todos = append(tl.Todos[:i], tl.Todos[i+1:]...)
			// Adjust selected index
			if tl.Selected >= len(tl.Todos) && tl.Selected > 0 {
				tl.Selected--
			}
			break
		}
	}

	return nil
}

// ToggleCompleted - Toggle status completed di database DAN memory sekaligus
func (tl *TodoList) ToggleCompleted(id int) error {
	// Get current todo
	var todo *Todo
	for _, t := range tl.Todos {
		if t.ID == id {
			todo = t
			break
		}
	}
	if todo == nil {
		return fmt.Errorf("todo with id %d not found", id)
	}

	newStatus := !todo.Completed
	query := `UPDATE todos SET completed=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`

	_, err := tl.db.Exec(query, newStatus, id)
	if err != nil {
		return fmt.Errorf("failed to toggle todo in database: %w", err)
	}

	// Update di memory
	todo.Completed = newStatus

	return nil
}

// GetByPriority - Filter (hanya memory)
func (tl *TodoList) GetByPriority(priority Priority) []*Todo {
	var filtered []*Todo
	for _, todo := range tl.Todos {
		if todo.Priority == priority {
			filtered = append(filtered, todo)
		}
	}
	return filtered
}

// GetCompleted - Filter (hanya memory)
func (tl *TodoList) GetCompleted() []*Todo {
	var completed []*Todo
	for _, todo := range tl.Todos {
		if todo.Completed {
			completed = append(completed, todo)
		}
	}
	return completed
}

// GetPending - Filter (hanya memory)
func (tl *TodoList) GetPending() []*Todo {
	var pending []*Todo
	for _, todo := range tl.Todos {
		if !todo.Completed {
			pending = append(pending, todo)
		}
	}
	return pending
}
