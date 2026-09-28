package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *DBConn
var OwnerDB *DBConn

// DBConn is a small wrapper around sqlx.DB that rebids `?` placeholders to
// Postgres-style `$1` placeholders automatically before executing queries.
type DBConn struct {
	db *sqlx.DB
}

func (c *DBConn) Query(query string, args ...interface{}) (*sql.Rows, error) {
	q := sqlx.Rebind(sqlx.DOLLAR, query)
	return c.db.Query(q, args...)
}

func (c *DBConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	q := sqlx.Rebind(sqlx.DOLLAR, query)
	return c.db.Exec(q, args...)
}

func (c *DBConn) QueryRow(query string, args ...interface{}) *sql.Row {
	q := sqlx.Rebind(sqlx.DOLLAR, query)
	return c.db.QueryRow(q, args...)
}

func Init() {
	var err error

	databaseURL := os.Getenv("DATABASE_URL")
	ownerDatabaseURL := os.Getenv("OWNER_DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}
	if ownerDatabaseURL == "" {
		ownerDatabaseURL = databaseURL
	}

	sqlDB, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Could not open main database: ", err)
	}
	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("Could not connect to main database: %v\nTroubleshooting hint: On Render, ensure both the Web Service and PostgreSQL database are in the SAME region if using the Internal Database URL, or use the External Database URL with sslmode=require.", err)
	}
	DB = &DBConn{db: sqlx.NewDb(sqlDB, "postgres")}

	sqlOwnerDB, err := sql.Open("postgres", ownerDatabaseURL)
	if err != nil {
		log.Fatal("Could not open owner database: ", err)
	}
	if err = sqlOwnerDB.Ping(); err != nil {
		log.Fatalf("Could not connect to owner database: %v\nTroubleshooting hint: On Render, ensure both the Web Service and PostgreSQL database are in the SAME region if using the Internal Database URL, or use the External Database URL with sslmode=require.", err)
	}
	OwnerDB = &DBConn{db: sqlx.NewDb(sqlOwnerDB, "postgres")}

	initDatabase(DB)
	initDatabase(OwnerDB)

	log.Println("Both databases connected and ready")
}

func initDatabase(conn *DBConn) {
	queries := []string{

		// Shops Table
		`CREATE TABLE IF NOT EXISTS shops (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		code TEXT UNIQUE NOT NULL,
		created_at TEXT DEFAULT now()::text
	);`,

		// Vendors Table
		`CREATE TABLE IF NOT EXISTS vendors (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		role TEXT NOT NULL DEFAULT 'vendor',
		shop_id TEXT NOT NULL,
		created_at TEXT DEFAULT now()::text,
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
	created_at TEXT DEFAULT now()::text,
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
		created_at TEXT DEFAULT now()::text,
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
		updated_at TEXT DEFAULT now()::text,
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
		created_at TEXT DEFAULT now()::text,
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
		created_at TEXT DEFAULT now()::text,
		FOREIGN KEY (vendor_id) REFERENCES users(id),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
	);`,
	}

	for _, q := range queries {
		if _, err := conn.Exec(q); err != nil {
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

func ensureColumn(conn *DBConn, table, column, definition string) error {
	// Information schema check for Postgres
	var colName string
	q := `SELECT column_name FROM information_schema.columns WHERE table_name = $1 AND column_name = $2`
	if err := conn.QueryRow(q, table, column).Scan(&colName); err == nil {
		return nil
	}
	// If not exists, add the column
	_, err := conn.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, definition))
	return err
}

func DBForRole(role string) *DBConn {
	if role == "owner" {
		return OwnerDB
	}
	return DB
}

func DBForEmail(email string) *DBConn {
	var u struct{ ID string }
	if err := OwnerDB.QueryRow(`SELECT id FROM users WHERE email = ?`, email).Scan(&u.ID); err == nil {
		return OwnerDB
	}
	return DB
}
