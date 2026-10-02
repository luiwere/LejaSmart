package test

import (
	"testing"

	"LejaSmart/db"
	"LejaSmart/models"
)

func TestAddSale(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	err := db.AddSale("vendor", "shop-1", "vendor-1", "Sugar", 10.0, 50.0, 30.0, "2024-01-15", "test sale")
	if err != nil {
		t.Fatalf("AddSale failed: %v", err)
	}

	sales, err := db.GetSales("vendor", "shop-1", "")
	if err != nil {
		t.Fatalf("GetSales failed: %v", err)
	}
	if len(sales) != 1 {
		t.Fatalf("expected 1 sale, got %d", len(sales))
	}
	if sales[0].ItemName != "Sugar" {
		t.Errorf("expected item_name 'Sugar', got '%s'", sales[0].ItemName)
	}
	if sales[0].Quantity != 10.0 {
		t.Errorf("expected quantity 10.0, got %f", sales[0].Quantity)
	}
	if sales[0].UnitPrice != 50.0 {
		t.Errorf("expected unit_price 50.0, got %f", sales[0].UnitPrice)
	}
	if sales[0].UnitCost != 30.0 {
		t.Errorf("expected unit_cost 30.0, got %f", sales[0].UnitCost)
	}
}

func TestGetSalesByShop(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	createTestUser(t, "vendor-2", "v2", "v2@test.com", "vendor", "shop-1")

	db.AddSale("vendor", "shop-1", "vendor-1", "Item A", 5.0, 100.0, 50.0, "2024-01-15", "")
	db.AddSale("vendor", "shop-1", "vendor-2", "Item B", 3.0, 200.0, 100.0, "2024-01-16", "")

	sales, err := db.GetSales("vendor", "shop-1", "")
	if err != nil {
		t.Fatalf("GetSales failed: %v", err)
	}
	if len(sales) != 2 {
		t.Fatalf("expected 2 sales, got %d", len(sales))
	}
}

func TestGetSalesByVendor(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	createTestUser(t, "vendor-2", "v2", "v2@test.com", "vendor", "shop-1")

	db.AddSale("vendor", "shop-1", "vendor-1", "Item A", 5.0, 100.0, 50.0, "2024-01-15", "")
	db.AddSale("vendor", "shop-1", "vendor-1", "Item B", 2.0, 150.0, 75.0, "2024-01-16", "")
	db.AddSale("vendor", "shop-1", "vendor-2", "Item C", 3.0, 200.0, 100.0, "2024-01-17", "")

	sales, err := db.GetSales("vendor", "shop-1", "vendor-1")
	if err != nil {
		t.Fatalf("GetSales failed: %v", err)
	}
	if len(sales) != 2 {
		t.Fatalf("expected 2 sales for vendor-1, got %d", len(sales))
	}
	for _, s := range sales {
		if s.VendorID != "vendor-1" {
			t.Errorf("expected vendor_id 'vendor-1', got '%s'", s.VendorID)
		}
	}
}

func TestGetSalesEmpty(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")

	sales, err := db.GetSales("vendor", "shop-1", "")
	if err != nil {
		t.Fatalf("GetSales failed: %v", err)
	}
	if len(sales) != 0 {
		t.Errorf("expected 0 sales, got %d", len(sales))
	}
}

func TestDeleteSale(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	err := db.AddSale("vendor", "shop-1", "vendor-1", "Sugar", 10.0, 50.0, 30.0, "2024-01-15", "")
	if err != nil {
		t.Fatalf("AddSale failed: %v", err)
	}

	sales, _ := db.GetSales("vendor", "shop-1", "")
	if len(sales) != 1 {
		t.Fatalf("expected 1 sale before delete, got %d", len(sales))
	}

	err = db.DeleteSale("vendor", "shop-1", sales[0].ID)
	if err != nil {
		t.Fatalf("DeleteSale failed: %v", err)
	}

	sales, _ = db.GetSales("vendor", "shop-1", "")
	if len(sales) != 0 {
		t.Errorf("expected 0 sales after delete, got %d", len(sales))
	}
}

func TestDeleteSaleNotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")

	err := db.DeleteSale("vendor", "shop-1", "nonexistent-id")
	if err != nil {
		t.Errorf("DeleteSale for non-existent ID should not error: %v", err)
	}
}

func TestGetSalesOrderDescending(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	db.AddSale("vendor", "shop-1", "vendor-1", "Item A", 1.0, 10.0, 5.0, "2024-01-10", "")
	db.AddSale("vendor", "shop-1", "vendor-1", "Item B", 1.0, 20.0, 10.0, "2024-01-20", "")
	db.AddSale("vendor", "shop-1", "vendor-1", "Item C", 1.0, 30.0, 15.0, "2024-01-15", "")

	sales, err := db.GetSales("vendor", "shop-1", "")
	if err != nil {
		t.Fatalf("GetSales failed: %v", err)
	}
	if len(sales) != 3 {
		t.Fatalf("expected 3 sales, got %d", len(sales))
	}
	if sales[0].Date != "2024-01-20" {
		t.Errorf("expected first sale date '2024-01-20', got '%s'", sales[0].Date)
	}
}

func TestSaleModel(t *testing.T) {
	s := models.Sale{
		ID:        "sale-1",
		VendorID:  "vendor-1",
		ItemName:  "Test Item",
		Quantity:  10.0,
		UnitPrice: 50.0,
		UnitCost:  30.0,
		Date:      "2024-01-15",
		Notes:     "test notes",
		CreatedAt: "2024-01-15T10:00:00Z",
	}
	if s.ItemName != "Test Item" {
		t.Errorf("expected item_name 'Test Item', got '%s'", s.ItemName)
	}
}
