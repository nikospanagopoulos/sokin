# Sokin

A personal progress journal and job-search assistant that runs locally on my own computer.
Built in Go, one small phase at a time, as a way to learn by building something I use every day.

> **Status: Phase 0 (planning).** Nothing is built yet. This README is the design, and it will be updated as each phase is finished.

## What it will do

- Every evening I write what I did today and what I plan for tomorrow, as a list of tasks.
- Each task has a progress value from 0 to 100%. The progress of a day is the average of its tasks. Days with no tasks are ignored.
- Weekly and monthly reviews show how much I achieved, what I learned, and which mistakes I repeated.
- A local AI model (Ollama) reads my tasks and my Obsidian notes and writes the reviews. Its suggestions are stored separately from what I wrote, so I can compare them.
- A "search now" button looks for job listings and scores how well they match my skills.
- Unfinished tasks are never moved automatically. The app asks me where each one should go.

## Roadmap

- [x] **Phase 0:** Plan, name, architecture, README
- [ ] **Phase 1:** Core journal: localhost web app, entries saved in PostgreSQL
- [ ] **Phase 2:** Tasks with progress and daily percentage
- [ ] **Phase 3:** Weekly and monthly reviews (plain calculations, no AI yet)
- [ ] **Phase 4:** Local AI reflection with Ollama and Obsidian notes
- [ ] **Phase 5:** Job search: listings, match score, application tracker
- [ ] **Phase 6:** Agent loop: ask in plain language, the app calls its own tools
- [ ] **Phase 7:** Protection: login, encryption, security testing on my own app
- [ ] **Phase 8:** Cloud: Docker, CI, deploy, backups, Telegram reminders

## Tech stack

| Part | Choice |
|---|---|
| Language | Go (standard library first) |
| Database | PostgreSQL, running in Docker |
| Frontend | HTML templates, plain CSS, a little JavaScript |
| AI | Ollama (local model) with a simple RAG approach |
| Notes | Obsidian vault (markdown files with frontmatter) |

## Architecture

One Go program split into packages (a modular monolith). Each phase adds a package without rewriting the old ones.

```
sokin/
├── cmd/sokin/        # entry point
├── internal/
│   ├── web/          # handlers and routes
│   ├── journal/      # daily entries
│   ├── tasks/        # tasks and percentages
│   ├── review/       # weekly and monthly reviews
│   └── store/        # PostgreSQL access
├── migrations/       # SQL files for the tables
├── templates/        # HTML
├── static/           # CSS and JS
└── docs/             # decisions and diagrams
    └── adr/          # architecture decision records
```

Three layers, one job each: `web` shows pages, the logic packages hold the rules, `store` is the only code that talks to the database. The percentage calculation knows nothing about HTTP or SQL, so it can be tested alone.

## Design principles

- **Percentages are calculated, not stored.** Only tasks are saved. Daily, weekly and monthly numbers come from the same data.
- **Works without JavaScript first.** The forms work with plain HTML. JavaScript only improves them. All calculations stay in Go.
- **Catch-up on startup.** The app does not run 24/7. When it opens, it checks for finished weeks or months without a review and creates them.
- **Read-only job search.** It only reads listings, never writes or sends anything, and only runs when I press the button. It respects each site's terms (APIs or RSS where possible).
- **Local first.** My notes stay on my computer. Cloud comes only in Phase 8, after security is in place.
- **Least privilege.** Every part gets only the access it needs.

## How to run

Coming in Phase 1.

## Privacy

This repository contains code and fake sample data only.
Real notes, `.env` files, tokens and database dumps are listed in `.gitignore` and are never committed.

## Why "Sokin"

"Sokin" is "Nikos" written backwards, my own name. The app works like a mirror of my real progress: it shows me what I did, what I learned and where I made mistakes.

It is a personal project by a student at Zone01 Athens. It tracks my progress, helps me find my first developer job, and teaches me the daily tools of Cloud/DevOps, AI and cybersecurity engineers.

## License

Sokin is licensed under the GNU Affero General Public License v3.0.
See the [LICENSE](LICENSE) file for details.