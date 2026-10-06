package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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

// NewDBConn wraps an existing *sql.DB into a DBConn. This is primarily used
// for testing where the underlying database may be SQLite instead of Postgres.
func NewDBConn(sqlDB *sql.DB) *DBConn {
	return &DBConn{db: sqlx.NewDb(sqlDB, "postgres")}
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
    	shop_id TEXT,
    	phone TEXT DEFAULT '',
    	address TEXT DEFAULT '',
    	bio TEXT DEFAULT '',
    	avatar_url TEXT DEFAULT '',
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

	// Owners are not attached to a shop, so users.shop_id has to accept NULL.
	if err := makeUserShopIDOptional(conn); err != nil {
		log.Fatal("could not relax users.shop_id:", err)
	}

	additionalColumns := map[string]map[string]string{
		"inventory": {
			"supplier_name": "supplier_name TEXT",
			"status":        "status TEXT",
			"reorder_level": "reorder_level REAL",
			"expiry_date":   "expiry_date TEXT",
			"restocked_at":  "restocked_at TEXT",
		},
		"users": {
			"phone":      "phone TEXT DEFAULT ''",
			"address":    "address TEXT DEFAULT ''",
			"bio":        "bio TEXT DEFAULT ''",
			"avatar_url": "avatar_url TEXT DEFAULT ''",
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

const userShopIDForeignKey = "users_shop_id_fkey"

// makeUserShopIDOptional lets users rows carry a NULL shop_id, which is how
// owners are stored. It is idempotent, so it is safe to run on every boot.
func makeUserShopIDOptional(conn *DBConn) error {
	var exists bool
	if err := conn.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM information_schema.table_constraints
		 WHERE table_name = 'users' AND constraint_name = $1)`,
		userShopIDForeignKey,
	).Scan(&exists); err != nil {
		return err
	}

	// An empty string can never satisfy the foreign key, so clear it first.
	if _, err := conn.Exec(`UPDATE users SET shop_id = NULL WHERE shop_id = ''`); err != nil {
		return err
	}

	if exists {
		if _, err := conn.Exec(fmt.Sprintf(`ALTER TABLE users DROP CONSTRAINT %s`, userShopIDForeignKey)); err != nil {
			return err
		}
	}

	// Re-create it without NOT NULL so NULL shop ids are accepted.
	if _, err := conn.Exec(`ALTER TABLE users ALTER COLUMN shop_id DROP NOT NULL`); err != nil {
		return err
	}
	_, err := conn.Exec(fmt.Sprintf(
		`ALTER TABLE users ADD CONSTRAINT %s FOREIGN KEY (shop_id) REFERENCES shops(id)`, userShopIDForeignKey))
	return err
}

func ensureColumn(conn *DBConn, table, column, definition string) error {
	// Information schema check for Postgres
	var colName string
	q := `SELECT column_name FROM information_schema.columns WHERE table_name = $1 AND column_name = $2`
	if err := conn.QueryRow(q, table, column).Scan(&colName); err == nil {
		return nil
	}
	// Some callers pass the full "column TYPE" fragment, so drop a leading
	// duplicate of the column name before building the statement.
	definition = strings.TrimSpace(definition)
	if rest := strings.TrimPrefix(definition, column); rest != definition {
		definition = strings.TrimSpace(rest)
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
