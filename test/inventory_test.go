package test

import (
	"testing"

	"LejaSmart/db"
	"LejaSmart/models"
)

func TestSaveInventoryItem(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	err := db.SaveInventoryItem(
		"vendor", "shop-1", "vendor-1", "Maize Flour", "Supplier A",
		"In Stock", "2024-12-31", "2024-01-01", 50.0, 10.0, "kg",
	)
	if err != nil {
		t.Fatalf("SaveInventoryItem failed: %v", err)
	}

	items, err := db.GetInventory("vendor", "shop-1", "vendor-1")
	if err != nil {
		t.Fatalf("GetInventory failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 inventory item, got %d", len(items))
	}
	if items[0].Name != "Maize Flour" {
		t.Errorf("expected name 'Maize Flour', got '%s'", items[0].Name)
	}
	if items[0].Quantity != 10.0 {
		t.Errorf("expected quantity 10.0, got %f", items[0].Quantity)
	}
	if items[0].Status != "In Stock" {
		t.Errorf("expected status 'In Stock', got '%s'", items[0].Status)
	}
}

func TestGetInventoryByVendor(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	createTestUser(t, "vendor-2", "v2", "v2@test.com", "vendor", "shop-1")

	db.SaveInventoryItem("vendor", "shop-1", "vendor-1", "Item A", "Sup A", "In Stock", "2024-12-31", "2024-01-01", 5.0, 100.0, "kg")
	db.SaveInventoryItem("vendor", "shop-1", "vendor-1", "Item B", "Sup B", "Low Stock", "2024-12-31", "2024-01-02", 3.0, 50.0, "bags")
	db.SaveInventoryItem("vendor", "shop-1", "vendor-2", "Item C", "Sup C", "In Stock", "2024-12-31", "2024-01-03", 10.0, 25.0, "units")

	items, err := db.GetInventory("vendor", "shop-1", "vendor-1")
	if err != nil {
		t.Fatalf("GetInventory failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items for vendor-1, got %d", len(items))
	}
}

func TestGetInventoryEmpty(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	items, err := db.GetInventory("vendor", "shop-1", "vendor-1")
	if err != nil {
		t.Fatalf("GetInventory failed: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestGetInventoryNotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	items, err := db.GetInventory("vendor", "nonexistent-shop", "nonexistent-vendor")
	if err != nil {
		t.Fatalf("GetInventory should not fail for non-existent data: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestSaveInventoryItemMultiple(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	items := []struct {
		name        string
		quantity    float64
		reorderLevel float64
	}{
		{"Sugar", 25.0, 5.0},
		{"Salt", 10.0, 3.0},
		{"Flour", 50.0, 10.0},
	}

	for _, item := range items {
		err := db.SaveInventoryItem(
			"vendor", "shop-1", "vendor-1", item.name, "Supplier", "In Stock",
			"2024-12-31", "2024-01-01", item.reorderLevel, item.quantity, "kg",
		)
		if err != nil {
			t.Fatalf("SaveInventoryItem failed for %s: %v", item.name, err)
		}
	}

	result, err := db.GetInventory("vendor", "shop-1", "vendor-1")
	if err != nil {
		t.Fatalf("GetInventory failed: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 items, got %d", len(result))
	}
}

func TestInventoryItemModel(t *testing.T) {
	i := models.InventoryItem{
		ID:           "inv-1",
		VendorID:     "vendor-1",
		Name:         "Test Item",
		SupplierName: "Supplier A",
		Status:       "In Stock",
		ReorderLevel: 5.0,
		ExpiryDate:   "2024-12-31",
		RestockedAt:  "2024-01-01",
		Quantity:     10.0,
		Unit:         "kg",
		UpdatedAt:    "2024-01-15T10:00:00Z",
	}
	if i.Name != "Test Item" {
		t.Errorf("expected name 'Test Item', got '%s'", i.Name)
	}
}
