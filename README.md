# 📸 Sabiflow

**Sabiflow** is a modern, photography-focused CRM application designed to help photographers manage their business workflows efficiently. It streamlines client management, job tracking, quotes, invoices, and scheduling—all within a self-hosted system. Built with Go (backend), SvelteKit (frontend), PostgreSQL (DB), and S3-compatible object storage.

---

## 🚀 Tech Stack

- **Backend:** Go (modular internal structure)
- **Frontend:** SvelteKit (SPA)
- **Database:** PostgreSQL (ORM/abstraction friendly)
- **Storage:** S3 (AWS or compatible)
- **Web Server:** Caddy (reverse proxy + TLS)
- **Migrations:** Goose
- **Email:** SMTP (templated email support)

---

## 📁 Project Structure

```
sabiflow/
├── backend/
│   ├── cmd/app/                  # main.go
│   ├── internal/
│   │   ├── handlers/             # HTTP API handlers
│   │   ├── services/             # Business logic
│   │   ├── repositories/         # DB access (sqlc or ent)
│   │   ├── models/               # Domain types and validation
│   │   ├── storage/              # S3 integration
│   │   ├── email/                # Email templating/sending
│   │   └── app/                  # Application bootstrap/deps
│   ├── migrations/               # Goose-compatible SQL files
│   ├── ui/                       # Built SvelteKit output (optional)
│   ├── Dockerfile
│   ├── Caddyfile
│   └── .env
└── frontend/
    ├── src/
    ├── static/
    ├── svelte.config.js
    └── vite.config.js
```

---

## 🎯 Milestone Plan (No Dates)

### 🔹 Milestone 1: Project Setup & Scaffolding
- Set up backend and frontend folder structure
- Add Dockerfile, .env, Makefile, Caddyfile
- Implement a health check endpoint and render a test frontend page

---

### 🔹 Milestone 2: Core Data Models & DB Layer
- Create schema for Clients, Jobs, Quotes, Tasks
- Add migrations
- Scaffold repository layer (sqlc or ent)
- Integrate DB into `Application` struct

---

### 🔹 Milestone 3: Lead Capture & Client Management
- Build public lead form
- Create `POST /api/leads` endpoint
- Develop admin UI for client CRUD operations
- Add notes and statuses

---

### 🔹 Milestone 4: Job Management & Task Scheduling
- Implement job creation and stages
- Add task scheduling (shoot/edit due dates)
- Introduce simple internal reminders

---

### 🔹 Milestone 5: Quotes & Invoices
- Develop quote builder UI
- Enable PDF generation (gofpdf or wkhtmltopdf)
- Implement email delivery
- Track payments manually

---

### 🔹 Milestone 6: Email Templates & Automation
- Create admin email template editor
- Support placeholders (e.g., `{{ClientName}}`)
- Enable manual and basic automated email triggers

---

### 🔹 Milestone 7: Polishing & Readiness
- Add CSRF protection, validation, and authentication
- Implement admin user management
- Introduce audit logs and basic logging
- Finalize UI/UX polish

---

### 🌞 Roadmap (Future)
- Contract builder (PDF + e-sign)
- Workflow pipelines
- Calendar sync (Google/Outlook)
- Multi-user/team support
- Client portal

---

## 📜 License

This project is licensed under the GNU General Public License v3.0.

You are free to use, modify, and distribute this software under the terms of the GPLv3. For full license details, see the [LICENSE](LICENSE) file.

