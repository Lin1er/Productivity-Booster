# ProdBooster 🚀

> A terminal-based productivity app for people who live in the terminal and struggle to stay organized.

[![Built with Go](https://img.shields.io/badge/built%20with-Go-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![TUI](https://img.shields.io/badge/interface-TUI-blueviolet?style=flat-square)](https://github.com/charmbracelet/bubbletea)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)

## Why I Built This 💭

Let's be real - I was a mess when it came to productivity.

Every time I opened my laptop or desktop, I'd forget what I was supposed to do. I'd have a million browser tabs open, sticky notes scattered everywhere (both physical and digital), and absolutely zero system to keep track of anything. I'd start the day with good intentions, then end up browsing Reddit for 3 hours wondering where the time went.

I tried Notion, Todoist, Google Keep - they're all great apps, but here's the thing: **I live in the terminal**. Switching to a browser or GUI app breaks my flow. I needed something that's always there, right in my face, the moment I log in.

So I built ProdBooster - a simple, no-nonsense TUI (Terminal User Interface) app that forces me to plan my day before I can do anything else. It's like having a personal accountability buddy that lives in your terminal.

## Features ✨

### The Three Pillars

**📋 Tasks (Todos)**

- Create tasks with priorities (High/Medium/Low)
- Set deadlines and get visual warnings for overdue items
- Color-coded by urgency (red = overdue, yellow = today, green = done)
- Mark tasks complete with a single spacebar press

**📅 Calendar (Events)**

- Schedule events with start/end times
- Add locations (perfect for meeting rooms or Zoom links)
- See what's happening today vs upcoming
- Past events fade out automatically

**📝 Notes**

- Quick capture for random thoughts and ideas
- Auto-timestamped so you know when you wrote it
- Newest notes appear first
- Perfect for meeting notes, brainstorming, or just dumping your brain

### Quality of Life Features

- **🔍 Smart Search & Filters** - Find anything instantly across all your data
- **🎨 Color Coding** - Visual cues so you know what needs attention at a glance
- **⌨️ Keyboard-First** - Everything is just a keystroke away, no mouse needed
- **📊 Dashboard** - See your day at a glance with 3 focused cards
- **💾 SQLite Storage** - Your data persists between sessions (obviously)
- **🎯 Auto-Launch Ready** - Perfect for TTY auto-launch on login (my original use case!)

## Screenshots 📸

_Dashboard - Your productivity command center_

```
┌────────────────────────────────────────────────────────┐
│ ✨ Wednesday, December 25, 2025 • 5 tasks pending     │
│    • 2 overdue • 1 due today                          │
├──────────┬──────────────┬──────────────┬──────────────┤
│ ✅ Tasks │ 📅 Events    │ 📒 Notes     │              │
│          │              │              │              │
│ ● Fix    │ 📅 Meeting   │ 📝 Ideas     │              │
│ ◐ Review │ 🕐 2:00 PM   │ 💭 Brain     │              │
│ ○ Update │              │    dump      │              │
└──────────┴──────────────┴──────────────┴──────────────┘
```

## Installation 🛠️

### Prerequisites

- Go 1.21 or higher
- SQLite3 (usually pre-installed on Linux/macOS)

### Quick Start

```bash
# Clone the repo
git clone https://github.com/Lin1er/Productivity-Booster.git
cd Productivity-Booster

# Build it
go build -o prodbooster

# Run it
./prodbooster
```

### Auto-Launch on TTY Login (Arch Linux)

This is how I use it - the app launches automatically when I log into my TTY (no display manager):

1. Add to your `.zprofile` or `.bash_profile`:

```bash
if [[ -z $DISPLAY ]] && [[ $(tty) = /dev/tty1 ]]; then
    ~/path/to/prodbooster
fi
```

2. Now when you log in to TTY1, you're forced to see your tasks before doing anything else. No more forgetting what you need to do!

## Usage Guide 🎮

### Navigation

- `Tab` - Switch between pages (Dashboard → Tasks → Calendar → Notes)
- `↑/↓` - Browse through lists
- `q` - Quit the app

### Tasks Page

- `n` - Create new task
- `e` - Edit selected task
- `d` - Delete selected task
- `Space` - Mark task as done/undone
- `/` - Search & filter

### Calendar Page

- `n` - Create new event
- `e` - Edit selected event
- `d` - Delete selected event
- `/` - Search & filter

### Notes Page

- `n` - Create new note
- `e` - Edit selected note
- `d` - Delete selected note
- `/` - Search & filter

### Dashboard

- `Tab` - Switch focus between cards
- `a` - Quick add (creates item in focused card)
- `Enter` - Jump to the focused page

## The Stack 🔧

Built with these awesome libraries:

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - The TUI framework that makes this possible
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - Pre-built TUI components (lists, inputs, etc.)
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** - Styling and layout (CSS for the terminal!)
- **[SQLite](https://www.sqlite.org/)** - Lightweight database for data persistence

## Project Structure 📁

```
.
├── main.go                 # Entry point
├── internal/
│   ├── db/                 # Database layer
│   │   └── db.go
│   ├── models/             # Data models (Todo, Note, Event)
│   │   ├── todo.go
│   │   ├── note.go
│   │   ├── event.go
│   │   └── navigation.go
│   └── ui/                 # User interface
│       ├── components/     # Reusable UI components
│       │   ├── todoForm.go
│       │   ├── noteForm.go
│       │   ├── eventForm.go
│       │   ├── searchBar.go
│       │   └── topbar.go
│       ├── pages/          # Full page views
│       │   ├── dashboard.go
│       │   ├── todos.go
│       │   ├── notes.go
│       │   └── calendar.go
│       └── styles/         # Global styles
│           └── main.go
└── tools/                  # Development tools
    └── cmd_seed.go         # Database seeder
```

## Development 👨‍💻

### Seeding Sample Data

Want to test with some sample data?

```bash
go run tools/cmd_seed.go
```

This creates:

- 5 sample todos (with varying priorities and due dates)
- 15 sample events (past, today, and future)
- 20 sample notes

### Building

```bash
# Development build
go build -o prodbooster

# Production build (smaller binary)
go build -ldflags="-s -w" -o prodbooster
```

## Roadmap 🗺️

Things I might add (or you can contribute!):

- [ ] Recurring tasks/events
- [ ] Tags for better organization
- [ ] Export to markdown/CSV
- [ ] Reminders/notifications (maybe using system notifications)
- [ ] Pomodoro timer integration
- [ ] Daily/weekly reports
- [ ] Sync across devices (maybe via git?)
- [ ] Vim-style keybindings option
- [ ] Custom themes

## Contributing 🤝

Found a bug? Have an idea? PRs are welcome!

1. Fork it
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License 📄

MIT License - do whatever you want with it!

## Acknowledgments 🙏

- Huge thanks to the [Charm](https://charm.sh/) team for creating Bubble Tea and the entire ecosystem
- Shoutout to everyone who struggles with productivity - we're in this together!
- Coffee. Lots of coffee.

## Contact 📬

If you want to chat about productivity, terminal apps, or just say hi:

- GitHub: [@Lin1er](https://github.com/Lin1er)
- Project: [Productivity-Booster](https://github.com/Lin1er/Productivity-Booster)

---

**Built with ❤️ and frustration with traditional todo apps**

_"The best productivity system is the one you actually use." - Me, probably_
