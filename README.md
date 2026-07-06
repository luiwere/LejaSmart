# DigiLedger

> A web-based expense tracking and business management platform built for local vendors who operate on paper. DigiLedger digitizes the financial record-keeping process — expenses, inventory, sales, and profit & loss — with role-based access for vendors, accountants, and business owners.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the App](#running-the-app)
- [User Roles](#user-roles)
- [API Endpoints](#api-endpoints)
- [Database Schema](#database-schema)
- [Authentication & Sessions](#authentication--sessions)
- [Security](#security)
- [Known Limitations](#known-limitations)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

DigiLedger was built to solve a real problem — local vendors and small business owners in informal markets still rely on paper ledgers to track their finances. This leads to lost records, calculation errors, and no easy way for an accountant to review a vendor's books remotely.

DigiLedger provides:

- A simple web interface that works on any device (phone, tablet, or desktop) with no installation required
- Role-based access so vendors log their own data while accountants get a birds-eye view
- An owner role for sole proprietors who handle both vendor and accountant responsibilities
- Voice input for hands-free expense logging — useful for vendors who are busy serving customers
- Real-time profit and loss calculation across any date range

---

## Features

### Vendor
- Log expenses with amount, date, category, supplier name, and notes
- Record expenses via voice input (Web Speech API)
- Track current inventory/supply levels with item name, quantity, and unit
- View profit and loss summary with optional date range filter
- Delete individual expense entries

### Accountant
- View all vendors' expenses in one place
- Filter expenses by individual vendor
- View total expense summary across all vendors
- Add new vendors to the system

### Owner
- Full access to both vendor and accountant features in a single tabbed dashboard
- Completely separate database from vendor/accountant data
- Manages their own shop independently

### General
- Secure registration and login with bcrypt password hashing
- Cookie-based session management with 24-hour expiry
- Back-button protection after logout (bfcache prevention)
- No-cache headers on protected pages
- Responsive layout — works on mobile and desktop

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go (standard library) |
| Frontend | HTML, CSS, Vanilla JavaScript |
| Database | SQLite (via `go-sqlite3`) |
| Password Hashing | bcrypt (`golang.org/x/crypto/bcrypt`) |
| Unique IDs | UUID (`github.com/google/uuid`) |
| Voice Input | Web Speech API (browser native) |

No frontend frameworks or external CSS libraries are used — everything is built from scratch.

---

## Project Structure

```
DigiLedger/
│
├── main.go                        # Entry point — server setup, routes, middleware
├── go.mod                         # Go module definition
├── go.sum                         # Dependency checksums (auto-generated)
├── README.md                      # Project documentation
├── .gitignore                     # Ignored files (*.db, .env, binaries)
│
├── db/
│   ├── database.go                # DB connection, table creation (two databases)
│   ├── users.go                   # User registration and lookup queries
│   ├── vendors.go                 # Vendor CRUD queries
│   ├── expenses.go                # Expense CRUD queries
│   ├── inventory.go               # Inventory CRUD queries
│   └── pnl.go                     # Profit & loss calculation queries
│
├── models/
│   ├── user.go                    # User struct
│   ├── vendor.go                  # Vendor struct
│   ├── expense.go                 # Expense struct
│   ├── inventory.go               # InventoryItem struct
│   └── pnl.go                     # PnLSummary struct
│
├── handlers/
│   ├── auth.go                    # Login, register, logout, /me endpoint
│   ├── vendors.go                 # Vendor and accountant dashboard handlers
│   ├── expenses.go                # Expense API handler
│   ├── inventory.go               # Inventory API handler
│   └── pnl.go                     # Profit & loss API handler
│
├── static/
│   ├── css/
│   │   └── style.css              # All styling — layout, components, responsive
│   └── js/
│       └── app.js                 # Client-side logic — fetch calls, voice input, tables
│
└── templates/
    ├── login.html                 # Login page
    ├── register.html              # Registration page
    ├── vendor-dashboard.html      # Vendor view
    ├── accountant-dashboard.html  # Accountant view
    └── owner-dashboard.html       # Combined owner view (tabbed)
```

---

## Getting Started

### Prerequisites

Make sure you have the following installed:

- [Go 1.21+](https://golang.org/dl/)
- GCC (required by `go-sqlite3` for CGo compilation)
  - **Ubuntu/Debian:** `sudo apt install gcc`
  - **Mac:** `xcode-select --install`
  - **Windows:** Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
- [Git](https://git-scm.com/)

### Installation

**1. Clone the repository**

```bash
git clone https://github.com/yourusername/DigiLedger.git
cd DigiLedger
```

**2. Install Go dependencies**

```bash
go mod tidy
```

This will install:
- `github.com/mattn/go-sqlite3`
- `github.com/google/uuid`
- `golang.org/x/crypto/bcrypt`

**3. Verify your project compiles**

```bash
go build ./...
```

If no errors appear, you are ready to run.

### Running the App

```bash
go run main.go
```

The server starts on port `8080`. Open your browser and visit:

```
http://localhost:8080
```

The SQLite database files (`digiledger.db` and `digiledger_owner.db`) are created automatically on first run — you do not need to set them up manually.

---

## User Roles

DigiLedger has three distinct roles, each with a different level of access:

| Role | Description | Dashboard |
|---|---|---|
| `vendor` | Logs expenses, inventory, and views their own P&L | `/vendor` |
| `accountant` | Views all vendors' data, manages vendor accounts | `/accountant` |
| `owner` | Does both — manages their own expenses AND views all vendor data | `/owner` |

### How roles work at registration

- Select your role from the dropdown when registering
- Vendors and accountants share the **same database** — they work together in the same shop
- Owners get their **own separate database** — their records are completely independent

### Vendor + Accountant relationship

A vendor and accountant who work at the same shop share data automatically. The accountant can see all of the vendor's logged expenses without any extra setup.

---

## API Endpoints

### Authentication

| Method | Route | Description | Auth required |
|---|---|---|---|
| `GET` | `/` | Login page | No |
| `POST` | `/login` | Submit login credentials | No |
| `GET` | `/register` | Registration page | No |
| `POST` | `/register` | Submit registration form | No |
| `GET` | `/logout` | Clear session and redirect to login | Yes |
| `GET` | `/me` | Returns logged-in user's ID and role (JSON) | Yes |

### Pages

| Method | Route | Description | Role |
|---|---|---|---|
| `GET` | `/vendor` | Vendor dashboard | vendor |
| `GET` | `/accountant` | Accountant dashboard | accountant |
| `GET` | `/owner` | Owner dashboard | owner |

### Data API

| Method | Route | Description |
|---|---|---|
| `GET` | `/expenses?vendorID=` | Get expenses (all or filtered by vendor) |
| `POST` | `/expenses` | Add a new expense |
| `DELETE` | `/expenses/{id}` | Delete an expense |
| `GET` | `/inventory?vendorID=` | Get inventory items |
| `POST` | `/inventory` | Add or update an inventory item |
| `GET` | `/pnl/{vendorID}` | Get P&L summary |
| `GET` | `/pnl/{vendorID}?from=&to=` | Get P&L for a date range |
| `GET` | `/vendors` | Get all vendors (accountant only) |
| `POST` | `/vendors` | Create a new vendor |

All data API routes return JSON.

---

## Database Schema

DigiLedger uses two SQLite databases with identical schemas:

- `digiledger.db` — used by vendors and accountants
- `digiledger_owner.db` — used exclusively by owners

### users

| Column | Type | Notes |
|---|---|---|
| id | TEXT | UUID primary key |
| username | TEXT | Display name |
| email | TEXT | Unique, used for login |
| password | TEXT | bcrypt hashed |
| role | TEXT | `vendor`, `accountant`, or `owner` |
| created_at | TEXT | Auto timestamp |

### vendors

| Column | Type | Notes |
|---|---|---|
| id | TEXT | UUID primary key |
| name | TEXT | Vendor display name |
| email | TEXT | Unique |
| role | TEXT | Default `vendor` |
| created_at | TEXT | Auto timestamp |

### expenses

| Column | Type | Notes |
|---|---|---|
| id | TEXT | UUID primary key |
| vendor_id | TEXT | Foreign key → vendors.id |
| amount | REAL | In KES |
| date | TEXT | ISO date string |
| category | TEXT | food, transport, supplies, utilities, other |
| supplier_name | TEXT | Optional |
| notes | TEXT | Optional |
| created_at | TEXT | Auto timestamp |

### inventory

| Column | Type | Notes |
|---|---|---|
| id | TEXT | UUID primary key |
| vendor_id | TEXT | Foreign key → vendors.id |
| name | TEXT | Item name |
| quantity | REAL | Numeric quantity |
| unit | TEXT | kg, litres, bags, etc. |
| updated_at | TEXT | Auto timestamp |

### income

| Column | Type | Notes |
|---|---|---|
| id | TEXT | UUID primary key |
| vendor_id | TEXT | Foreign key → vendors.id |
| amount | REAL | In KES |
| date | TEXT | ISO date string |
| notes | TEXT | Optional |
| created_at | TEXT | Auto timestamp |

---

## Authentication & Sessions

DigiLedger uses **cookie-based session management**:

1. On successful login, two `HttpOnly` cookies are set:
   - `session_user` — stores the logged-in user's UUID
   - `session_role` — stores the user's role (`vendor`, `accountant`, `owner`)

2. Both cookies expire after **24 hours**

3. On logout, both cookies are cleared by setting their expiry to the past

4. Protected routes (`/vendor`, `/accountant`, `/owner`) check for a valid `session_user` cookie before serving the page — unauthenticated requests are redirected to `/`

5. The `/me` endpoint lets the frontend read the current user's ID and role from the session without exposing it in localStorage

---

## Security

| Concern | How it's handled |
|---|---|
| Password storage | bcrypt hashing with `DefaultCost` — never stored as plain text |
| Duplicate email detection | SQLite `UNIQUE constraint` + `sqlite3.ErrConstraintUnique` error code check (not string matching) |
| Session cookies | `HttpOnly` flag prevents JavaScript access — protects against XSS |
| Back button after logout | `NoCacheMiddleware` sets `no-store, no-cache` headers + `pageshow` bfcache event listener in JS |
| Unauthenticated access | Session cookie check on all protected routes — redirects to login if missing |
| Owner data isolation | Separate SQLite database (`digiledger_owner.db`) — owner data never mixes with vendor/accountant data |

---

## Known Limitations

- **No income tracking UI yet** — the `income` table exists in the database but there is no form for vendors to log sales/revenue, which means P&L currently shows expenses only
- **No shop pairing system** — vendors and accountants are assumed to share the same shop but there is no formal shop ID or invite code system yet
- **No password reset** — users cannot recover a forgotten password
- **Sessions are not invalidated server-side on logout** — only the client cookie is cleared; if someone copied the cookie before logout, it would remain valid until expiry
- **Voice input is browser-dependent** — Chrome supports it well; Firefox and Safari have limited or no support for the Web Speech API

---

## Roadmap

- [ ] Sales/income tracking section with per-item revenue logging
- [ ] Automatic P&L calculation from sales minus expenses
- [ ] Shop ID system — vendor generates a code, accountant joins using it
- [ ] Password reset via email
- [ ] Server-side session invalidation on logout
- [ ] Export reports to PDF or Excel
- [ ] Docker deployment setup
- [ ] Unit tests for all handlers and db functions

## Docker

Run the app with Docker (multi-stage build with SQLite support):

1. Create a data directory for the SQLite files (persisted on host):

```bash
mkdir -p data
touch data/Digiledgerledger.db data/Digiledgerowner.db
```

2. Build the Docker image:

```bash
docker build -t digiledger:latest .
```

3. Run with Docker (bind mounts persist DB files):

```bash
docker run --rm -p 8080:8080 \
  -v "$PWD/data/Digiledgerledger.db":/app/Digiledgerledger.db \
  -v "$PWD/data/Digiledgerowner.db":/app/Digiledgerowner.db \
  digiledger:latest
```

Or use `docker compose` (included):

```bash
docker compose up --build
```

Notes:
- The image includes the system SQLite library so `go-sqlite3` works at runtime.
- If you prefer not to build locally, you can run inside the container via the provided `docker-compose.yml`.


---

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/your-feature-name`
3. Commit your changes using conventional commits: `git commit -m "feat: add income tracking form"`
4. Push to your branch: `git push origin feat/your-feature-name`
5. Open a pull request

Please make sure your code compiles (`go build ./...`) and follows Go formatting standards (`go fmt ./...`) before submitting.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

> Built with Go · SQLite · HTML · CSS · Vanilla JS

