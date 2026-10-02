package test

import (
	"testing"

	"LejaSmart/db"
	"LejaSmart/models"
)

func TestCreateVendor(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")

	err := db.CreateVendor("vendor", "shop-1", "New Vendor", "newvendor@test.com", "vendor")
	if err != nil {
		t.Fatalf("CreateVendor failed: %v", err)
	}

	vendors, err := db.GetAllVendors("vendor", "shop-1")
	if err != nil {
		t.Fatalf("GetAllVendors failed: %v", err)
	}
	if len(vendors) != 1 {
		t.Fatalf("expected 1 vendor, got %d", len(vendors))
	}
	if vendors[0].Name != "New Vendor" {
		t.Errorf("expected name 'New Vendor', got '%s'", vendors[0].Name)
	}
	if vendors[0].Email != "newvendor@test.com" {
		t.Errorf("expected email 'newvendor@test.com', got '%s'", vendors[0].Email)
	}
	if vendors[0].ShopID != "shop-1" {
		t.Errorf("expected shop_id 'shop-1', got '%s'", vendors[0].ShopID)
	}
}

func TestCreateVendorDuplicateEmail(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")

	err := db.CreateVendor("vendor", "shop-1", "Vendor A", "dup@test.com", "vendor")
	if err != nil {
		t.Fatalf("first CreateVendor failed: %v", err)
	}

	err = db.CreateVendor("vendor", "shop-1", "Vendor B", "dup@test.com", "vendor")
	if err == nil {
		t.Error("expected error for duplicate email")
	}
}

func TestGetAllVendors(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestShop(t, "shop-2", "Other Shop", "CODE2")

	db.CreateVendor("vendor", "shop-1", "Vendor A", "a@test.com", "vendor")
	db.CreateVendor("vendor", "shop-1", "Vendor B", "b@test.com", "vendor")

	vendors, err := db.GetAllVendors("vendor", "shop-1")
	if err != nil {
		t.Fatalf("GetAllVendors failed: %v", err)
	}
	if len(vendors) != 2 {
		t.Fatalf("expected 2 vendors, got %d", len(vendors))
	}

	shop2Vendors, err := db.GetAllVendors("vendor", "shop-2")
	if err != nil {
		t.Fatalf("GetAllVendors failed: %v", err)
	}
	if len(shop2Vendors) != 0 {
		t.Errorf("expected 0 vendors for shop-2, got %d", len(shop2Vendors))
	}
}

func TestGetAllVendorsEmpty(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")

	vendors, err := db.GetAllVendors("vendor", "shop-1")
	if err != nil {
		t.Fatalf("GetAllVendors failed: %v", err)
	}
	if len(vendors) != 0 {
		t.Errorf("expected 0 vendors, got %d", len(vendors))
	}
}

func TestGetAllVendorsFields(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")

	db.CreateVendor("vendor", "shop-1", "Test Vendor", "tv@test.com", "vendor")

	vendors, _ := db.GetAllVendors("vendor", "shop-1")
	if len(vendors) != 1 {
		t.Fatalf("expected 1 vendor, got %d", len(vendors))
	}
	v := vendors[0]
	if v.ID == "" {
		t.Error("expected non-empty ID")
	}
	if v.Role != "vendor" {
		t.Errorf("expected role 'vendor', got '%s'", v.Role)
	}
	if v.CreatedAt == "" {
		t.Error("expected non-empty created_at")
	}
}

func TestVendorModel(t *testing.T) {
	v := models.Vendor{
		ID:        "vendor-1",
		Name:      "Test Vendor",
		Email:     "test@test.com",
		Role:      "vendor",
		ShopID:    "shop-1",
		CreatedAt: "2024-01-15T10:00:00Z",
	}
	if v.Name != "Test Vendor" {
		t.Errorf("expected name 'Test Vendor', got '%s'", v.Name)
	}
}
