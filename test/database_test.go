package test

import (
	"database/sql"
	"testing"

	"LejaSmart/db"
)

func TestNewDBConn(t *testing.T) {
	sqlDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	defer sqlDB.Close()

	conn := db.NewDBConn(sqlDB)
	if conn == nil {
		t.Fatal("NewDBConn returned nil")
	}
}

func TestDBConnExec(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	res, err := db.DB.Exec(`INSERT INTO shops (id, name, code) VALUES (?, ?, ?)`,
		"shop-1", "Test Shop", "CODE1")
	if err != nil {
		t.Fatalf("DBConn.Exec failed: %v", err)
	}
	rows, _ := res.RowsAffected()
	if rows != 1 {
		t.Errorf("expected 1 row affected, got %d", rows)
	}
}

func TestDBConnQuery(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.DB.Exec(`INSERT INTO shops (id, name, code) VALUES (?, ?, ?)`,
		"shop-1", "Test Shop", "CODE1")
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	rows, err := db.DB.Query(`SELECT id, name, code FROM shops WHERE id = ?`, "shop-1")
	if err != nil {
		t.Fatalf("DBConn.Query failed: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, name, code string
		if err := rows.Scan(&id, &name, &code); err != nil {
			t.Fatalf("row scan failed: %v", err)
		}
		if id != "shop-1" {
			t.Errorf("expected id 'shop-1', got '%s'", id)
		}
		if name != "Test Shop" {
			t.Errorf("expected name 'Test Shop', got '%s'", name)
		}
		if code != "CODE1" {
			t.Errorf("expected code 'CODE1', got '%s'", code)
		}
		count++
	}
	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
	}
}

func TestDBConnQueryRow(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.DB.Exec(`INSERT INTO shops (id, name, code) VALUES (?, ?, ?)`,
		"shop-1", "Test Shop", "CODE1")
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	var name string
	err = db.DB.QueryRow(`SELECT name FROM shops WHERE id = ?`, "shop-1").Scan(&name)
	if err != nil {
		t.Fatalf("DBConn.QueryRow failed: %v", err)
	}
	if name != "Test Shop" {
		t.Errorf("expected 'Test Shop', got '%s'", name)
	}
}

func TestDBConnQueryRowNotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	var name string
	err := db.DB.QueryRow(`SELECT name FROM shops WHERE id = ?`, "nonexistent").Scan(&name)
	if err == nil {
		t.Error("expected error for non-existent row, got nil")
	}
}

func TestDBForRoleVendor(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	conn := db.DBForRole("vendor")
	if conn != db.DB {
		t.Error("DBForRole('vendor') should return db.DB")
	}
}

func TestDBForRoleOwner(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	conn := db.DBForRole("owner")
	if conn != db.OwnerDB {
		t.Error("DBForRole('owner') should return db.OwnerDB")
	}
}

func TestDBForRoleEmpty(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	conn := db.DBForRole("")
	if conn != db.DB {
		t.Error("DBForRole('') should return db.DB")
	}
}

func TestDBForEmailFoundInOwner(t *testing.T) {
	cleanup := setupTestDBs(t)
	defer cleanup()

	_, err := db.OwnerDB.Exec(
		`INSERT INTO users (id, username, email, password, role, shop_id) VALUES (?, ?, ?, ?, ?, ?)`,
		"owner-1", "owner_user", "owner@test.com", "hashed", "owner", nil,
	)
	if err != nil {
		t.Fatalf("failed to insert owner: %v", err)
	}

	conn := db.DBForEmail("owner@test.com")
	if conn != db.OwnerDB {
		t.Error("DBForEmail should return OwnerDB when email is found there")
	}
}

func TestDBForEmailFoundInMain(t *testing.T) {
	cleanup := setupTestDBs(t)
	defer cleanup()

	_, err := db.DB.Exec(
		`INSERT INTO users (id, username, email, password, role, shop_id) VALUES (?, ?, ?, ?, ?, ?)`,
		"vendor-1", "vendor_user", "vendor@test.com", "hashed", "vendor", "shop-1",
	)
	if err != nil {
		t.Fatalf("failed to insert vendor: %v", err)
	}

	conn := db.DBForEmail("vendor@test.com")
	if conn != db.DB {
		t.Error("DBForEmail should return DB when email is not found in OwnerDB")
	}
}

func TestDBForEmailNotFound(t *testing.T) {
	cleanup := setupTestDBs(t)
	defer cleanup()

	conn := db.DBForEmail("nonexistent@test.com")
	if conn != db.DB {
		t.Error("DBForEmail should return DB when email is not found in any database")
	}
}
