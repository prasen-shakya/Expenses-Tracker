# Expense Tracker

This is a small full-stack expense tracker built with Go, Postgres, and a plain HTML/CSS/JavaScript frontend.

The goal was to keep it simple: one backend, one database, one UI, and no extra moving parts unless they earned their place. The Go app serves both the API and the frontend, so once the server is running you can open the app and start logging expenses right away.

## What it does

- Lets users register and log in
- Stores passwords securely with bcrypt
- Uses JWT auth for protected expense routes
- Tracks vendor, amount, category, description, and date for each expense
- Shows a dashboard view with monthly totals, averages, recent activity, and category breakdowns
- Supports creating, editing, listing, and deleting expenses

## Stack

- Go for the backend
- PostgreSQL for persistence
- Vanilla HTML, CSS, and JavaScript for the frontend
- Vercel-compatible server entrypoint in `api/index.go`

## Project structure

```text
.
├── api/                # Vercel handler entrypoint
├── internal/           # App, auth, controllers, db, routes, domain logic
├── sql/                # Database schema
├── web/                # Embedded frontend assets
├── main.go             # Local server entrypoint
└── vercel.json         # Rewrites for deployment
```

## Running locally

### 1. Set up environment variables

Create a `.env` file in the project root with:

```env
DATABASE_URL=postgres://USERNAME:PASSWORD@HOST:5432/DB_NAME?sslmode=disable
JWT_SECRET=replace-this-with-a-real-secret
PORT=3000
```

Only `DATABASE_URL` and `JWT_SECRET` are required. If `PORT` is missing, the app defaults to `3001`.

### 2. Create the database tables

Run the schema file against your Postgres database:

```bash
psql "$DATABASE_URL" -f sql/schema.sql
```

### 3. Start the app

```bash
go run .
```

Then open [http://localhost:3001](http://localhost:3001).

## API routes

Public routes:

- `POST /register`
- `POST /login`

Protected routes:

- `GET /expenses`
- `GET /expenses/total`
- `GET /expenses/{id}`
- `POST /expenses`
- `PUT /expenses/{id}`
- `DELETE /expenses/{id}`

Protected routes expect a bearer token in the `Authorization` header:

```http
Authorization: Bearer <jwt-token>
```

## Example request

```bash
curl -X POST http://localhost:3001/expenses \
  -H "Authorization: Bearer <jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "vendor": "Trader Joe'\''s",
    "amount": 42.18,
    "category": "Groceries",
    "description": "Weekly groceries",
    "date": "2026-03-24T00:00:00Z"
  }'
```

## A couple of implementation notes

- The frontend is embedded into the Go binary using `embed`, so the server ships the UI directly.
- Session data is stored in `localStorage` on the frontend.
- Passwords must be at least 8 characters.
- If an expense date is missing or invalid, the server falls back to the current UTC time.

## Deploying

This repo includes a Vercel entrypoint at `api/index.go` and rewrite rules in `vercel.json`, so the same app can handle both the frontend and API routes in deployment.

You will need to provide the same environment variables in your deployed environment:

- `DATABASE_URL`
- `JWT_SECRET`

## Why this project exists

Mostly to keep expense tracking straightforward. A lot of personal tools either feel bloated or make you fight the UI for small tasks. This one is meant to feel lightweight: log the expense, check the month, move on.
