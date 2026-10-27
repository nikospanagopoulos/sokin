# 0001: Use PostgreSQL

Status: Accepted
Date: 2026-10-27

## Context

Sokin stores daily entries, tasks with a progress value (0-100%), and generated reviews. Tasks belong to a day, and the daily, weekly and monthly percentages are calculated from them, never stored. The app will keep asking questions such as "all tasks of this week" or "the average progress of last month".

I am also learning backend and DevOps skills, and I want to learn how a real SQL database server works as part of this project.

## Decision

Use PostgreSQL, running in Docker on my own computer. The table structure lives in SQL migration files in the repository. The real data stays in a Docker volume outside the repository.

## Alternatives considered

- **JSON or markdown files:** simple to start, but every weekly or monthly calculation would need my own code to load, filter and combine the data. That is what a database already does well.
- **SQLite:** also SQL, and much easier to set up, since it is one file. I rejected it because I want to learn a database that runs as a separate server, which is closer to what is used in backend and DevOps work, and it fits the Docker practice planned for the cloud phase.

## Consequences

**Good**
- Percentages and reviews become simple queries instead of custom code.
- I practice SQL, migrations, and running a service in Docker from the first phase.
- The same setup moves naturally to the cloud in Phase 8.

**Bad**
- The app cannot run without the database container, so setup is heavier than a single file.
- Docker on Windows (WSL2) uses a noticeable part of the 16 GB of RAM on my laptop.
- Backups are my responsibility. A backup contains personal data, so it must never go into the repository.