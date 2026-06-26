package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var OwnerDB *sql.DB

const (
	SharedDBPath = "./Digiledgerledger.db"
	OwnerDBPath  = "./Digiledgerowner.db"
)

func Init() {
	var err error

	// Open shared database for vendors and accountants
	DB, err = sql.Open("sqlite3", SharedDBPath)
	if err != nil {
		log.Fatal("Could not open shared Database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Could not connect to shared Database:", err)
	}
	initDatabase(DB)

	// Open owner-specific database
	OwnerDB, err = sql.Open("sqlite3", OwnerDBPath)
	if err != nil {
		log.Fatal("Could not open owner Database:", err)
	}
	if err = OwnerDB.Ping(); err != nil {
		log.Fatal("Could not connect to owner Database:", err)
	}
	initDatabase(OwnerDB)

	log.Println("Shared and owner databases connected and ready")
}

func initDatabase(conn *sql.DB) {
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
	email TEXT UNIQUE NOT NULL,
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

	// Income Table
	`CREATE TABLE IF NOT EXISTS income (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		notes TEXT,
		created_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (vendor_id) REFERENCES vendors(id)
	);`,

	// Sales Table
	`CREATE TABLE IF NOT EXISTS sales (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		item_name TEXT NOT NULL,
		quantity REAL NOT NULL,
		unit_price REAL NOT NULL,
		unit_cost REAL,
		date TEXT NOT NULL,
		notes TEXT,
		created_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (vendor_id) REFERENCES vendors(id)
	);`,
	}

	for _, q := range queries {
		_, err := conn.Exec(q)
		if err != nil {
			log.Fatal("could not create table:", err)
		}
	}
}

func DBForRole(role string) *sql.DB {
	if role == "owner" {
		return OwnerDB
	}
	return DB
}

func DBForEmail(email string) *sql.DB {
	var u struct{ ID string }
	if err := OwnerDB.QueryRow(`SELECT id FROM users WHERE email = ?`, email).Scan(&u.ID); err == nil {
		return OwnerDB
	}
	return DB
}
