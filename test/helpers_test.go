package test

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"LejaSmart/db"

	_ "github.com/mattn/go-sqlite3"
)

var testSchema = []string{
	`CREATE TABLE shops (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		code TEXT UNIQUE NOT NULL,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP
	);`,
	`CREATE TABLE vendors (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		role TEXT NOT NULL DEFAULT 'vendor',
		shop_id TEXT NOT NULL,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,
	`CREATE TABLE users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'vendor',
		shop_id TEXT,
		phone TEXT DEFAULT '',
		address TEXT DEFAULT '',
		bio TEXT DEFAULT '',
		avatar_url TEXT DEFAULT '',
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,
	`CREATE TABLE expenses (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		shop_id TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		category TEXT,
		supplier_name TEXT,
		notes TEXT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vendor_id) REFERENCES users(id),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,
	`CREATE TABLE inventory (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		shop_id TEXT NOT NULL,
		name TEXT NOT NULL,
		supplier_name TEXT,
		status TEXT,
		reorder_level REAL,
		expiry_date TEXT,
		restocked_at TEXT,
		quantity REAL NOT NULL,
		unit TEXT,
		updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vendor_id) REFERENCES users(id),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,
	`CREATE TABLE income (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		shop_id TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		notes TEXT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vendor_id) REFERENCES users(id),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,
	`CREATE TABLE sales (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		shop_id TEXT NOT NULL,
		item_name TEXT NOT NULL,
		quantity REAL NOT NULL,
		unit_price REAL NOT NULL,
		unit_cost REAL,
		date TEXT NOT NULL,
		notes TEXT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vendor_id) REFERENCES users(id),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,
}

func initializeSchema(t *testing.T, sqlDB *sql.DB) {
	t.Helper()
	for _, q := range testSchema {
		if _, err := sqlDB.Exec(q); err != nil {
			t.Fatalf("failed to create schema: %v", err)
		}
	}
}

func setupTestDB(t *testing.T) func() {
	t.Helper()
	sqlDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	initializeSchema(t, sqlDB)

	conn := db.NewDBConn(sqlDB)
	db.DB = conn
	db.OwnerDB = conn

	return func() {
		sqlDB.Close()
	}
}

func setupTestDBs(t *testing.T) func() {
	t.Helper()
	sqlMain, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	sqlMain.SetMaxOpenConns(1)
	initializeSchema(t, sqlMain)

	sqlOwner, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	sqlOwner.SetMaxOpenConns(1)
	initializeSchema(t, sqlOwner)

	db.DB = db.NewDBConn(sqlMain)
	db.OwnerDB = db.NewDBConn(sqlOwner)

	return func() {
		sqlMain.Close()
		sqlOwner.Close()
	}
}

func newAuthenticatedRequest(method, target, userID, role string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)
	req.AddCookie(&http.Cookie{Name: "session_user", Value: userID})
	req.AddCookie(&http.Cookie{Name: "session_role", Value: role})
	return req
}

func createTestShop(t *testing.T, id, name, code string) {
	t.Helper()
	_, err := db.DB.Exec(`INSERT INTO shops (id, name, code) VALUES (?, ?, ?)`, id, name, code)
	if err != nil {
		t.Fatalf("failed to create test shop: %v", err)
	}
}

func createTestUser(t *testing.T, id, username, email, role, shopID string) {
	t.Helper()
	_, err := db.DB.Exec(
		`INSERT INTO users (id, username, email, password, role, shop_id) VALUES (?, ?, ?, ?, ?, ?)`,
		id, username, email, "hashedpassword", role, shopID,
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
}
