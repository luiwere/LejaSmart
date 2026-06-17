package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error

	// Create or Open the SQLite Database file
	DB, err = sql.Open("sqlite3", "./Digiledgerledger.db")
	if err != nil {
		log.Fatal("Could not open Database:", err)
	}

	// Test Connection

	if err = DB.Ping(); err != nil {
		log.Fatal("Could not connect to te Database:", err)
	}

	// Create the tables
	createTables()
	log.Println("Database connected and ready")

}

func createTables() {
	queries := []string{

	// Vendors Table
	`CREATE TABLE IF NOT EXISTS vendors (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		role TEXT NOT NULL DEFAULT 'vendor',
		created_at TEXT DEFAULT (datetime('now'))
	);`,

	// Users Table
	`CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	username TEXT NOT NULL,
	email TEXT NOT NULL,
	password TEXT NOT NULL,
	role TEXT NOT NULL DEFAULT 'vendor',
	created_at TEXT DEFAULT (datetime('now'))
	);`,

	// Expenses Table
	`CREATE TABLE IF NOT EXISTS expenses (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		category TEXT,
		supplier_name TEXT,
		notes TEXT,
		created_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (vendor_id) REFERENCES vendors(id)

	);`,

	// Inventory Table
	`CREATE TABLE IF NOT EXISTS inventory (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		name TEXT NOT NULL,
		quantity REAL NOT NULL,
		unit TEXT,
		updated_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (vendor_id) REFERENCES vendors(id)
	);`,

	// Income Timetable
	`CREATE TABLE IF NOT EXISTS income (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		notes TEXT,
		created_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (vendor_id) REFERENCES vendorS(id)

	);`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatal("coult not create Table:", err)
		}
	}
}
