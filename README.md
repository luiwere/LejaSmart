# LejaSmart

> A web-based expense tracking and business management platform built for local vendors who operate on paper. LejaSmart digitizes the financial record-keeping process — expenses, inventory, sales, and profit & loss — with role-based access for vendors, accountants, and business owners.

---

## Overview

LejaSmart helps local vendors and small business owners replace paper ledgers with a simple web app for tracking expenses, inventory, sales, and profit & loss. It supports three roles (vendor, accountant, owner) and uses SQLite for persistence.

## Project Structure (high level)

```
LejaSmart/
├── main.go
├── go.mod
├── README.md
├── db/
├── models/
├── handlers/
├── static/
└── templates/
```

## Getting Started (short)

Clone and run:

```bash
git clone https://github.com/yourusername/LejaSmart.git
cd LejaSmart
go mod tidy
go run main.go
```

Visit http://localhost:8080

The app creates two SQLite files automatically on first run: `lejasmart.db` and `lejasmart_owner.db`.

## Docker (short)

Create data dir and run with Docker:

```bash
mkdir -p data
touch data/lejasmart.db data/lejasmart_owner.db
docker build -t lejasmart:latest .
docker run --rm -p 8080:8080 \
  -v "$PWD/data/lejasmart.db":/app/lejasmart.db \
  -v "$PWD/data/lejasmart_owner.db":/app/lejasmart_owner.db \
  lejasmart:latest
```

For full documentation, refer to the other project files.
