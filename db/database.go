package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var OwnerDB *sql.DB

func Init() {
	var err error

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "."
	}

	// Open shared database for vendors and accountants
	DB, err = sql.Open("sqlite3", dbPath+"/lejasmart.db")
	if err != nil {
		log.Fatal("Could not open main database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Could not connect to main database:", err)
	}

	OwnerDB, err = sql.Open("sqlite3", dbPath+"/lejasmart_owner.db")
	if err != nil {
		log.Fatal("Could not open owner database:", err)
	}
	if err = OwnerDB.Ping(); err != nil {
		log.Fatal("Could not connect to owner database:", err)
	}

	initDatabase(DB)
	initDatabase(OwnerDB)

	log.Println("Both databases connected and ready")
}

func initDatabase(conn *sql.DB) {
	queries := []string{

		// Shops Table
		`CREATE TABLE IF NOT EXISTS shops (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		code TEXT UNIQUE NOT NULL,
		created_at TEXT DEFAULT (datetime('now'))
	);`,

		// Vendors Table
		`CREATE TABLE IF NOT EXISTS vendors (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		role TEXT NOT NULL DEFAULT 'vendor',
		shop_id TEXT NOT NULL,
		created_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,

		// Users Table
		`CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	username TEXT NOT NULL,
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	role TEXT NOT NULL DEFAULT 'vendor',
	shop_id TEXT NOT NULL,
	created_at TEXT DEFAULT (datetime('now')),
	FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,

		// Expenses Table
		`CREATE TABLE IF NOT EXISTS expenses (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		shop_id TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		category TEXT,
		supplier_name TEXT,
		notes TEXT,
		created_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (vendor_id) REFERENCES users(id),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,

		// Inventory Table
		`CREATE TABLE IF NOT EXISTS inventory (
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
		updated_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (vendor_id) REFERENCES users(id),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,

		// Income Table
		`CREATE TABLE IF NOT EXISTS income (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		shop_id TEXT NOT NULL,
		amount REAL NOT NULL,
		date TEXT NOT NULL,
		notes TEXT,
		created_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (vendor_id) REFERENCES users(id),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,

		// Sales Table
		`CREATE TABLE IF NOT EXISTS sales (
		id TEXT PRIMARY KEY,
		vendor_id TEXT NOT NULL,
		shop_id TEXT NOT NULL,
		item_name TEXT NOT NULL,
		quantity REAL NOT NULL,
		unit_price REAL NOT NULL,
		unit_cost REAL,
		date TEXT NOT NULL,
		notes TEXT,
		created_at TEXT DEFAULT (datetime('now')),
		FOREIGN KEY (vendor_id) REFERENCES users(id),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,
	}

	for _, q := range queries {
		_, err := conn.Exec(q)
		if err != nil {
			log.Fatal("could not create table:", err)
		}
	}

	columnsToEnsure := map[string]string{
		"users":     "shop_id TEXT NOT NULL DEFAULT ''",
		"vendors":   "shop_id TEXT NOT NULL DEFAULT ''",
		"expenses":  "shop_id TEXT NOT NULL DEFAULT ''",
		"inventory": "shop_id TEXT NOT NULL DEFAULT ''",
		"income":    "shop_id TEXT NOT NULL DEFAULT ''",
		"sales":     "shop_id TEXT NOT NULL DEFAULT ''",
	}

	for table, definition := range columnsToEnsure {
		if err := ensureColumn(conn, table, "shop_id", definition); err != nil {
			log.Fatal("could not ensure column:", err)
		}
	}

	additionalColumns := map[string]map[string]string{
		"inventory": {
			"supplier_name": "supplier_name TEXT",
			"status":        "status TEXT",
			"reorder_level": "reorder_level REAL",
			"expiry_date":   "expiry_date TEXT",
			"restocked_at":  "restocked_at TEXT",
		},
	}

	for table, columns := range additionalColumns {
		for column, definition := range columns {
			if err := ensureColumn(conn, table, column, definition); err != nil {
				log.Fatal("could not ensure column:", err)
			}
		}
	}
}

func ensureColumn(conn *sql.DB, table, column, definition string) error {
	rows, err := conn.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var colName string
		var colType string
		var notnull int
		var dflt sql.NullString
		var pk int
		if err := rows.Scan(&cid, &colName, &colType, &notnull, &dflt, &pk); err != nil {
			return err
		}
		if colName == column {
			return nil
		}
	}

	_, err = conn.Exec(`ALTER TABLE ` + table + ` ADD COLUMN ` + column + ` ` + definition)
	return err
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
