package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Note struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}

type NoteList struct {
	db       *sql.DB
	Notes    []*Note
	Selected int
	NextID   int
}

// Load - Load semua notes dari database ke memory
func (nl *NoteList) Load() error {
	query := "SELECT id, title, content, created_at FROM notes ORDER BY id"
	rows, err := nl.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	nl.Notes = []*Note{} // Clear existing

	for rows.Next() {
		var id int
		var title, content string
		var createdAt time.Time

		if err := rows.Scan(&id, &title, &content, &createdAt); err != nil {
			return fmt.Errorf("failed to scan note: %w", err)
		}

		note := &Note{
			ID:        id,
			Title:     title,
			Content:   content,
			CreatedAt: createdAt,
		}

		nl.Notes = append(nl.Notes, note)

		// Update NextID
		if id >= nl.NextID {
			nl.NextID = id + 1
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating notes: %w", err)
	}

	return nil
}

/*
* Sample NoteList for testing purpose
 */

func NewNoteListSample() *NoteList {
	return &NoteList{
		Notes: []*Note{
			{
				ID: 1, Title: "Meeting Notes",
				Content:   "Discuss project roadmap and milestones.",
				CreatedAt: time.Now(),
			},
			{
				ID: 2, Title: "Shopping List",
				Content: "Eggs, Milk, Bread, Butter.",
			},
			{
				ID:      3,
				Title:   "Ideas",
				Content: "Start a blog about Go programming.",
			},
		},
	}
}

func NewNoteList(db_ *sql.DB) *NoteList {
	nl := &NoteList{
		db:       db_,
		Notes:    []*Note{},
		Selected: 0,
		NextID:   1,
	}
	// Auto-load dari database saat inisialisasi
	if err := nl.Load(); err != nil {
		// Log error tapi tetap return instance kosong
		fmt.Printf("Warning: failed to load notes: %v\n", err)
	}
	return nl
}

func (nl *NoteList) SelectNext() {
	if nl.Selected < len(nl.Notes)-1 {
		nl.Selected++
	}
}

func (nl *NoteList) SelectPrev() {
	if nl.Selected > 0 {
		nl.Selected--
	}
}

func (nl *NoteList) GetSelected() *Note {
	if len(nl.Notes) == 0 || nl.Selected >= len(nl.Notes) {
		return nil
	}
	return nl.Notes[nl.Selected]
}

func (nl *NoteList) GetSelectedTitle() string {
	if len(nl.Notes) == 0 || nl.Selected >= len(nl.Notes) {
		return ""
	}
	return nl.Notes[nl.Selected].Title
}

func (nl *NoteList) GetSelectedContent() string {
	if len(nl.Notes) == 0 || nl.Selected >= len(nl.Notes) {
		return ""
	}
	return nl.Notes[nl.Selected].Content
}

func (nl NoteList) Count() int {
	return len(nl.Notes)
}

// Add - Tambah note ke database DAN memory sekaligus
func (nl *NoteList) Add(title, content string) error {
	query := `INSERT INTO notes (title, content, created_at) VALUES (?, ?, ?)`

	now := time.Now()
	result, err := nl.db.Exec(query, title, content, now)
	if err != nil {
		return fmt.Errorf("failed to add note to database: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Tambah ke memory
	note := &Note{
		ID:        int(id),
		Title:     title,
		Content:   content,
		CreatedAt: now,
	}
	nl.Notes = append(nl.Notes, note)
	nl.NextID = int(id) + 1

	return nil
}

// Update - Update note di database DAN memory sekaligus
func (nl *NoteList) Update(id int, title, content string) error {
	query := `UPDATE notes SET title=?, content=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`

	_, err := nl.db.Exec(query, title, content, id)
	if err != nil {
		return fmt.Errorf("failed to update note in database: %w", err)
	}

	// Update di memory
	for _, note := range nl.Notes {
		if note.ID == id {
			note.Title = title
			note.Content = content
			break
		}
	}

	return nil
}

// Remove - Hapus note dari database DAN memory sekaligus
func (nl *NoteList) Remove(id int) error {
	query := `DELETE FROM notes WHERE id=?`

	_, err := nl.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete note from database: %w", err)
	}

	// Hapus dari memory
	for i, note := range nl.Notes {
		if note.ID == id {
			nl.Notes = append(nl.Notes[:i], nl.Notes[i+1:]...)
			// Adjust selected index
			if nl.Selected >= len(nl.Notes) && nl.Selected > 0 {
				nl.Selected--
			}
			break
		}
	}

	return nil
}
