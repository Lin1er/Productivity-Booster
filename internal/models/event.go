package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Event struct {
	ID        int // kebutuhan database
	Title     string
	Content   string
	Location  string
	StartTime time.Time
	EndTime   time.Time
}

type EventList struct {
	db       *sql.DB
	Events   []*Event
	Selected int
	NextID   int
}

// Load - Load semua events dari database ke memory
func (el *EventList) Load() error {
	query := "SELECT id, title, description, location, start_time, end_time FROM events ORDER BY start_time"
	rows, err := el.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	el.Events = []*Event{} // Clear existing

	for rows.Next() {
		var id int
		var title, description, location string
		var startTime, endTime time.Time

		if err := rows.Scan(&id, &title, &description, &location, &startTime, &endTime); err != nil {
			return fmt.Errorf("failed to scan event: %w", err)
		}

		event := &Event{
			ID:        id,
			Title:     title,
			Content:   description,
			Location:  location,
			StartTime: startTime,
			EndTime:   endTime,
		}

		el.Events = append(el.Events, event)

		// Update NextID
		if id >= el.NextID {
			el.NextID = id + 1
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating events: %w", err)
	}

	return nil
}

func NewEventList(db_ *sql.DB) *EventList {
	el := &EventList{
		db:       db_,
		Events:   []*Event{},
		Selected: 0,
		NextID:   1,
	}
	// Auto-load dari database saat inisialisasi
	if err := el.Load(); err != nil {
		// Log error tapi tetap return instance kosong
		fmt.Printf("Warning: failed to load events: %v\n", err)
	}
	return el
}

func (el *EventList) SelectNext() {
	if el.Selected < len(el.Events)-1 {
		el.Selected++
	}
}

func (el *EventList) SelectPrev() {
	if el.Selected > 0 {
		el.Selected--
	}
}

func (el *EventList) Count() int {
	return len(el.Events)
}

func (el *EventList) GetSelected() *Event {
	if len(el.Events) == 0 || el.Selected >= len(el.Events) {
		return nil
	}
	return el.Events[el.Selected]
}

func (el *EventList) GetTodayEvents() []*Event {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)
	var todayEvents []*Event

	for _, event := range el.Events {
		if event.StartTime.After(today) && event.StartTime.Before(tomorrow) {
			todayEvents = append(todayEvents, event)
		}
	}
	return todayEvents
}

func (el *EventList) GetUpcomingEvents(count int) []*Event {
	now := time.Now()
	var upcomingEvents []*Event
	for _, event := range el.Events {
		if event.StartTime.After(now) {
			upcomingEvents = append(upcomingEvents, event)
		}
		if len(upcomingEvents) >= count {
			break
		}
	}
	return upcomingEvents
}

// Add - Tambah event ke database DAN memory sekaligus
func (el *EventList) Add(title, content, location string, startTime, endTime time.Time) error {
	query := `INSERT INTO events (title, description, location, start_time, end_time, created_at)
	          VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := el.db.Exec(query, title, content, location, startTime, endTime, now)
	if err != nil {
		return fmt.Errorf("failed to add event to database: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Tambah ke memory
	event := &Event{
		ID:        int(id),
		Title:     title,
		Content:   content,
		Location:  location,
		StartTime: startTime,
		EndTime:   endTime,
	}
	el.Events = append(el.Events, event)
	el.NextID = int(id) + 1

	return nil
}

// Update - Update event di database DAN memory sekaligus
func (el *EventList) Update(id int, title, content, location string, startTime, endTime time.Time) error {
	query := `UPDATE events SET title=?, description=?, location=?, start_time=?, end_time=?, updated_at=CURRENT_TIMESTAMP
	          WHERE id=?`

	_, err := el.db.Exec(query, title, content, location, startTime, endTime, id)
	if err != nil {
		return fmt.Errorf("failed to update event in database: %w", err)
	}

	// Update di memory
	for _, event := range el.Events {
		if event.ID == id {
			event.Title = title
			event.Content = content
			event.Location = location
			event.StartTime = startTime
			event.EndTime = endTime
			break
		}
	}

	return nil
}

// Remove - Hapus event dari database DAN memory sekaligus
func (el *EventList) Remove(id int) error {
	query := `DELETE FROM events WHERE id=?`

	_, err := el.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event from database: %w", err)
	}

	// Hapus dari memory
	for i, event := range el.Events {
		if event.ID == id {
			el.Events = append(el.Events[:i], el.Events[i+1:]...)
			// Adjust selected index
			if el.Selected >= len(el.Events) && el.Selected > 0 {
				el.Selected--
			}
			break
		}
	}

	return nil
}
