package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Get database path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
	}

	dbDir := filepath.Join(homeDir, ".prodbooster")
	dbPath := filepath.Join(dbDir, "data.db")

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	fmt.Println("🌱 Seeding database with sample data...")

	// Seed Todos (5 items)
	seedTodos(db)

	// Seed Events (15 items)
	seedEvents(db)

	// Seed Notes (20 items)
	seedNotes(db)

	fmt.Println("✅ Database seeded successfully!")
}

func seedTodos(db *sql.DB) {
	now := time.Now()
	todos := []struct {
		title       string
		description string
		completed   bool
		priority    int
		dueDate     *time.Time
	}{
		{
			title:       "Review project proposal",
			description: "Review and provide feedback on Q1 2026 project proposals",
			completed:   false,
			priority:    2, // High
			dueDate:     timePtr(now.Add(2 * time.Hour)),
		},
		{
			title:       "Team standup meeting",
			description: "Daily standup with development team at 10 AM",
			completed:   false,
			priority:    1, // Medium
			dueDate:     timePtr(now.Add(24 * time.Hour)),
		},
		{
			title:       "Fix critical bug in production",
			description: "Database connection timeout issue reported by users",
			completed:   false,
			priority:    2, // High
			dueDate:     timePtr(now.Add(-2 * time.Hour)), // Overdue!
		},
		{
			title:       "Update documentation",
			description: "Update API documentation for v2.0 release",
			completed:   false,
			priority:    0, // Low
			dueDate:     timePtr(now.Add(72 * time.Hour)),
		},
		{
			title:       "Code review for PR #234",
			description: "Review authentication refactoring pull request",
			completed:   true,
			priority:    1, // Medium
			dueDate:     timePtr(now.Add(-24 * time.Hour)),
		},
	}

	for _, todo := range todos {
		_, err := db.Exec(`
			INSERT INTO todos (title, description, completed, priority, created_at, due_date)
			VALUES (?, ?, ?, ?, ?, ?)
		`, todo.title, todo.description, todo.completed, todo.priority, now, todo.dueDate)

		if err != nil {
			log.Printf("Error seeding todo '%s': %v\n", todo.title, err)
		} else {
			fmt.Printf("✓ Added todo: %s\n", todo.title)
		}
	}
}

func seedEvents(db *sql.DB) {
	now := time.Now()
	events := []struct {
		title       string
		description string
		location    string
		startTime   time.Time
		endTime     time.Time
	}{
		// Today's events
		{
			title:       "Morning Standup",
			description: "Daily team sync meeting",
			location:    "Zoom",
			startTime:   time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location()),
			endTime:     time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location()),
		},
		{
			title:       "Client Demo",
			description: "Product demo for potential client",
			location:    "Meeting Room A",
			startTime:   time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, now.Location()),
			endTime:     time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, now.Location()),
		},
		{
			title:       "Gym Session",
			description: "Workout routine - Chest and triceps",
			location:    "Local Gym",
			startTime:   time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location()),
			endTime:     time.Date(now.Year(), now.Month(), now.Day(), 19, 30, 0, 0, now.Location()),
		},
		// Tomorrow's events
		{
			title:       "Sprint Planning",
			description: "Planning for Sprint 24",
			location:    "Conference Room B",
			startTime:   now.Add(24*time.Hour + 10*time.Hour),
			endTime:     now.Add(24*time.Hour + 12*time.Hour),
		},
		{
			title:       "1-on-1 with Manager",
			description: "Weekly sync with team lead",
			location:    "Zoom",
			startTime:   now.Add(24*time.Hour + 15*time.Hour),
			endTime:     now.Add(24*time.Hour + 16*time.Hour),
		},
		{
			title:       "Dinner with Sarah",
			description: "Birthday dinner celebration",
			location:    "Italian Restaurant Downtown",
			startTime:   now.Add(24*time.Hour + 19*time.Hour),
			endTime:     now.Add(24*time.Hour + 21*time.Hour),
		},
		// This week
		{
			title:       "Architecture Review",
			description: "System design review for microservices migration",
			location:    "Zoom",
			startTime:   now.Add(48*time.Hour + 13*time.Hour),
			endTime:     now.Add(48*time.Hour + 15*time.Hour),
		},
		{
			title:       "Coffee with John",
			description: "Catch up with old colleague",
			location:    "Starbucks Central",
			startTime:   now.Add(48*time.Hour + 16*time.Hour),
			endTime:     now.Add(48*time.Hour + 17*time.Hour),
		},
		{
			title:       "Product Roadmap Meeting",
			description: "Q1 2026 planning session",
			location:    "Conference Room A",
			startTime:   now.Add(72*time.Hour + 10*time.Hour),
			endTime:     now.Add(72*time.Hour + 12*time.Hour),
		},
		{
			title:       "Team Building Activity",
			description: "Monthly team outing",
			location:    "Bowling Alley",
			startTime:   now.Add(72*time.Hour + 17*time.Hour),
			endTime:     now.Add(72*time.Hour + 20*time.Hour),
		},
		// Next week
		{
			title:       "Conference: DevOps Summit 2026",
			description: "Annual DevOps conference - Day 1",
			location:    "Convention Center",
			startTime:   now.Add(168 * time.Hour),
			endTime:     now.Add(168*time.Hour + 8*time.Hour),
		},
		{
			title:       "Dentist Appointment",
			description: "Regular checkup",
			location:    "Dr. Smith Dental Clinic",
			startTime:   now.Add(168*time.Hour + 14*time.Hour),
			endTime:     now.Add(168*time.Hour + 15*time.Hour),
		},
		// Past events
		{
			title:       "Code Review Session",
			description: "Weekly code review meeting",
			location:    "Zoom",
			startTime:   now.Add(-24 * time.Hour),
			endTime:     now.Add(-23 * time.Hour),
		},
		{
			title:       "All Hands Meeting",
			description: "Company-wide quarterly update",
			location:    "Main Auditorium",
			startTime:   now.Add(-48 * time.Hour),
			endTime:     now.Add(-46 * time.Hour),
		},
		{
			title:       "Hackathon Weekend",
			description: "Internal innovation hackathon",
			location:    "Office",
			startTime:   now.Add(-72 * time.Hour),
			endTime:     now.Add(-60 * time.Hour),
		},
	}

	for _, event := range events {
		_, err := db.Exec(`
			INSERT INTO events (title, description, location, start_time, end_time)
			VALUES (?, ?, ?, ?, ?)
		`, event.title, event.description, event.location, event.startTime, event.endTime)

		if err != nil {
			log.Printf("Error seeding event '%s': %v\n", event.title, err)
		} else {
			fmt.Printf("✓ Added event: %s\n", event.title)
		}
	}
}

func seedNotes(db *sql.DB) {
	now := time.Now()
	notes := []struct {
		title   string
		content string
	}{
		{
			title:   "Project Ideas",
			content: "1. Personal finance tracker with AI insights\n2. Pomodoro timer with analytics\n3. Habit tracker with gamification\n4. Local-first note-taking app\n5. CLI tool for managing dotfiles",
		},
		{
			title:   "Book Notes: Clean Code",
			content: "Key takeaways:\n- Meaningful names are crucial\n- Functions should do one thing\n- Comments don't make up for bad code\n- Error handling is important\n- Write tests first (TDD)",
		},
		{
			title:   "Meeting Notes: Sprint Retrospective",
			content: "What went well:\n- Faster deployment pipeline\n- Better communication\n- Code quality improved\n\nWhat to improve:\n- Documentation needs work\n- More pair programming\n- Better estimation",
		},
		{
			title:   "Learning Go - Best Practices",
			content: "1. Use gofmt for formatting\n2. Handle errors explicitly\n3. Use defer for cleanup\n4. Prefer composition over inheritance\n5. Keep interfaces small\n6. Use context for cancellation",
		},
		{
			title:   "Recipe: Spaghetti Carbonara",
			content: "Ingredients:\n- 400g spaghetti\n- 200g guanciale or pancetta\n- 4 egg yolks\n- 100g Pecorino Romano\n- Black pepper\n\nCook pasta, fry guanciale, mix with eggs and cheese off heat. Simple!",
		},
		{
			title:   "Workout Plan - Week 1",
			content: "Monday: Chest & Triceps\nTuesday: Back & Biceps\nWednesday: Rest\nThursday: Legs\nFriday: Shoulders & Abs\nWeekend: Cardio + Stretching",
		},
		{
			title:   "Interview Questions - Backend Dev",
			content: "Technical:\n- Explain REST vs GraphQL\n- Database indexing strategies\n- Microservices vs Monolith\n- Caching strategies\n- Rate limiting implementation\n\nBehavioral:\n- Biggest technical challenge\n- Conflict resolution\n- Learning new technology",
		},
		{
			title:   "Travel Checklist",
			content: "Before trip:\n☐ Book flights and hotels\n☐ Check passport validity\n☐ Travel insurance\n☐ Notify bank\n☐ Download offline maps\n☐ Pack essentials\n☐ Set home security",
		},
		{
			title:   "Side Project: Task Manager",
			content: "Features to implement:\n- User authentication\n- Real-time sync\n- Tags and filters\n- Calendar view\n- Mobile responsive\n- Dark mode\n- Export/import data\n\nTech stack: Go + HTMX + SQLite",
		},
		{
			title:   "Daily Journal - Dec 17",
			content: "Today was productive! Finished the database integration for ProdBooster. The TUI is coming together nicely.\n\nLearned about Bubble Tea's message passing system - it's quite elegant.\n\nTomorrow: Focus on search functionality and UI polish.",
		},
		{
			title:   "Gift Ideas for Christmas",
			content: "Mom: Kitchen gadget set\nDad: Wireless headphones\nSister: Art supplies\nBrother: Gaming mouse\nGrandma: Photo album\nBest friend: Coffee subscription",
		},
		{
			title:   "Debugging Tips",
			content: "1. Read the error message carefully\n2. Check the logs first\n3. Reproduce the issue consistently\n4. Isolate the problem\n5. Use rubber duck debugging\n6. Take breaks when stuck\n7. Ask for help after 30min",
		},
		{
			title:   "Useful Terminal Commands",
			content: "# Find large files\ndu -ah / | sort -rh | head -n 20\n\n# Monitor system resources\nhtop\n\n# Network debugging\nss -tulpn\n\n# Search in files\nrg \"pattern\" --type go\n\n# Git aliases\ngit log --oneline --graph",
		},
		{
			title:   "Productivity Tips",
			content: "Morning routine:\n- No phone for first hour\n- Exercise or stretch\n- Healthy breakfast\n- Review daily goals\n\nWork habits:\n- Pomodoro technique (25min focus)\n- Batch similar tasks\n- Time blocking\n- Turn off notifications\n- Regular breaks",
		},
		{
			title:   "System Design Notes",
			content: "CAP Theorem:\n- Consistency\n- Availability\n- Partition Tolerance\n\nCan only guarantee 2 out of 3!\n\nScaling strategies:\n- Vertical: More powerful machines\n- Horizontal: More machines\n- Caching: Redis, Memcached\n- Load balancing: Round-robin, least connections",
		},
		{
			title:   "Go Idioms & Patterns",
			content: "// Check error immediately\nif err != nil {\n    return err\n}\n\n// Defer for cleanup\nf, _ := os.Open(file)\ndefer f.Close()\n\n// Accept interfaces, return structs\nfunc Process(r io.Reader) *Result\n\n// Use blank identifier\n_, err := doSomething()",
		},
		{
			title:   "Meal Prep Ideas",
			content: "Sunday prep:\n- Chicken breast (grilled)\n- Brown rice (batch cook)\n- Roasted vegetables\n- Hard boiled eggs\n- Overnight oats\n- Cut fruits\n\nEasy to mix and match throughout the week!",
		},
		{
			title:   "Learning Resources",
			content: "Go:\n- Effective Go (official)\n- Go by Example\n- Ardan Labs courses\n\nSystem Design:\n- System Design Primer (GitHub)\n- Designing Data-Intensive Applications\n\nAlgorithms:\n- LeetCode\n- Neetcode.io",
		},
		{
			title:   "Bug Tracker - ProdBooster",
			content: "Fixed:\n✓ Form keyboard navigation conflict\n✓ Search enter key not working\n\nIn Progress:\n- Dashboard stats calculation\n- Event color coding\n\nBacklog:\n- Error notifications\n- Input validation\n- Confirmation dialogs",
		},
		{
			title:   "Monthly Goals - December",
			content: "Personal:\n☐ Finish ProdBooster v1.0\n☐ Read 2 books\n☐ Exercise 4x/week\n☐ Meditate daily\n\nProfessional:\n☐ Complete certification course\n☐ Contribute to open source\n☐ Write 2 blog posts\n☐ Learn HTMX",
		},
	}

	for _, note := range notes {
		_, err := db.Exec(`
			INSERT INTO notes (title, content, created_at)
			VALUES (?, ?, ?)
		`, note.title, note.content, now)

		if err != nil {
			log.Printf("Error seeding note '%s': %v\n", note.title, err)
		} else {
			fmt.Printf("✓ Added note: %s\n", note.title)
		}
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
