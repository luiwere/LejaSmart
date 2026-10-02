package test

import (
	"testing"

	"LejaSmart/db"
	"LejaSmart/models"
)

func TestAddExpense(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "testvendor", "v@test.com", "vendor", "shop-1")

	err := db.AddExpense("vendor", "shop-1", "vendor-1", 250.50, "2024-01-15", "food", "Supplier A", "test notes")
	if err != nil {
		t.Fatalf("AddExpense failed: %v", err)
	}

	expenses, err := db.GetExpenses("vendor", "shop-1", "")
	if err != nil {
		t.Fatalf("GetExpenses failed: %v", err)
	}
	if len(expenses) != 1 {
		t.Fatalf("expected 1 expense, got %d", len(expenses))
	}
	if expenses[0].Amount != 250.50 {
		t.Errorf("expected amount 250.50, got %f", expenses[0].Amount)
	}
	if expenses[0].Category != "food" {
		t.Errorf("expected category 'food', got '%s'", expenses[0].Category)
	}
}

func TestGetExpensesByShop(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	createTestUser(t, "vendor-2", "v2", "v2@test.com", "vendor", "shop-1")

	db.AddExpense("vendor", "shop-1", "vendor-1", 100.0, "2024-01-15", "food", "Sup A", "notes")
	db.AddExpense("vendor", "shop-1", "vendor-2", 200.0, "2024-01-16", "transport", "Sup B", "")

	expenses, err := db.GetExpenses("vendor", "shop-1", "")
	if err != nil {
		t.Fatalf("GetExpenses failed: %v", err)
	}
	if len(expenses) != 2 {
		t.Fatalf("expected 2 expenses, got %d", len(expenses))
	}
}

func TestGetExpensesByVendor(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	createTestUser(t, "vendor-2", "v2", "v2@test.com", "vendor", "shop-1")

	db.AddExpense("vendor", "shop-1", "vendor-1", 100.0, "2024-01-15", "food", "Sup A", "")
	db.AddExpense("vendor", "shop-1", "vendor-1", 150.0, "2024-01-16", "transport", "Sup B", "")
	db.AddExpense("vendor", "shop-1", "vendor-2", 200.0, "2024-01-17", "food", "Sup C", "")

	expenses, err := db.GetExpenses("vendor", "shop-1", "vendor-1")
	if err != nil {
		t.Fatalf("GetExpenses failed: %v", err)
	}
	if len(expenses) != 2 {
		t.Fatalf("expected 2 expenses for vendor-1, got %d", len(expenses))
	}
	for _, e := range expenses {
		if e.VendorID != "vendor-1" {
			t.Errorf("expected vendor_id 'vendor-1', got '%s'", e.VendorID)
		}
	}
}

func TestGetExpensesEmpty(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")

	expenses, err := db.GetExpenses("vendor", "shop-1", "")
	if err != nil {
		t.Fatalf("GetExpenses failed: %v", err)
	}
	if len(expenses) != 0 {
		t.Errorf("expected 0 expenses, got %d", len(expenses))
	}
}

func TestDeleteExpense(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	err := db.AddExpense("vendor", "shop-1", "vendor-1", 100.0, "2024-01-15", "food", "Sup A", "")
	if err != nil {
		t.Fatalf("AddExpense failed: %v", err)
	}

	expenses, _ := db.GetExpenses("vendor", "shop-1", "")
	if len(expenses) != 1 {
		t.Fatalf("expected 1 expense before delete, got %d", len(expenses))
	}

	err = db.DeleteExpense("vendor", "shop-1", expenses[0].ID)
	if err != nil {
		t.Fatalf("DeleteExpense failed: %v", err)
	}

	expenses, _ = db.GetExpenses("vendor", "shop-1", "")
	if len(expenses) != 0 {
		t.Errorf("expected 0 expenses after delete, got %d", len(expenses))
	}
}

func TestDeleteExpenseNotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")

	err := db.DeleteExpense("vendor", "shop-1", "nonexistent-id")
	if err != nil {
		t.Errorf("DeleteExpense for non-existent ID should not error: %v", err)
	}
}

func TestExpensesOrderDescending(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	db.AddExpense("vendor", "shop-1", "vendor-1", 100.0, "2024-01-10", "food", "Sup A", "")
	db.AddExpense("vendor", "shop-1", "vendor-1", 200.0, "2024-01-20", "food", "Sup B", "")
	db.AddExpense("vendor", "shop-1", "vendor-1", 50.0, "2024-01-15", "food", "Sup C", "")

	expenses, err := db.GetExpenses("vendor", "shop-1", "")
	if err != nil {
		t.Fatalf("GetExpenses failed: %v", err)
	}
	if len(expenses) != 3 {
		t.Fatalf("expected 3 expenses, got %d", len(expenses))
	}
	if expenses[0].Date != "2024-01-20" {
		t.Errorf("expected first expense date '2024-01-20', got '%s'", expenses[0].Date)
	}
}

func TestExpenseModelFields(t *testing.T) {
	e := models.Expense{
		ID:           "exp-1",
		VendorID:     "vendor-1",
		Amount:       99.99,
		Date:         "2024-01-15",
		Category:     "food",
		SupplierName: "Supplier A",
		Notes:        "test notes",
		CreatedAt:    "2024-01-15T10:00:00Z",
	}
	if e.Amount != 99.99 {
		t.Errorf("expected amount 99.99, got %f", e.Amount)
	}
}
