package test

import (
	"encoding/json"
	"testing"

	"LejaSmart/models"
)

func TestExpenseJSON(t *testing.T) {
	original := models.Expense{
		ID:           "exp-1",
		VendorID:     "vendor-1",
		Amount:       99.99,
		Date:         "2024-01-15",
		Category:     "food",
		SupplierName: "Supplier A",
		Notes:        "test notes",
		CreatedAt:    "2024-01-15T10:00:00Z",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if !containsKey(data, "id") {
		t.Error("JSON should contain 'id' key")
	}
	if !containsKey(data, "vendor_id") {
		t.Error("JSON should contain 'vendor_id' key")
	}
	if !containsKey(data, "amount") {
		t.Error("JSON should contain 'amount' key")
	}
	if !containsKey(data, "category") {
		t.Error("JSON should contain 'category' key")
	}
	if !containsKey(data, "supplier_name") {
		t.Error("JSON should contain 'supplier_name' key")
	}

	var decoded models.Expense
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, original.ID)
	}
	if decoded.Amount != original.Amount {
		t.Errorf("Amount mismatch: got %f, want %f", decoded.Amount, original.Amount)
	}
	if decoded.Category != original.Category {
		t.Errorf("Category mismatch: got %s, want %s", decoded.Category, original.Category)
	}
}

func TestInventoryItemJSON(t *testing.T) {
	original := models.InventoryItem{
		ID:           "inv-1",
		VendorID:     "vendor-1",
		Name:         "Maize Flour",
		SupplierName: "Supplier A",
		Status:       "In Stock",
		ReorderLevel: 5.0,
		ExpiryDate:   "2024-12-31",
		RestockedAt:  "2024-01-01",
		Quantity:     10.0,
		Unit:         "kg",
		UpdatedAt:    "2024-01-15T10:00:00Z",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if !containsKey(data, "reorder_level") {
		t.Error("JSON should contain 'reorder_level' key")
	}
	if !containsKey(data, "expiry_date") {
		t.Error("JSON should contain 'expiry_date' key")
	}
	if !containsKey(data, "restocked_at") {
		t.Error("JSON should contain 'restocked_at' key")
	}

	var decoded models.InventoryItem
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, original.Name)
	}
	if decoded.ReorderLevel != original.ReorderLevel {
		t.Errorf("ReorderLevel mismatch: got %f, want %f", decoded.ReorderLevel, original.ReorderLevel)
	}
}

func TestSaleJSON(t *testing.T) {
	original := models.Sale{
		ID:        "sale-1",
		VendorID:  "vendor-1",
		ItemName:  "Sugar",
		Quantity:  10.0,
		UnitPrice: 50.0,
		UnitCost:  30.0,
		Date:      "2024-01-15",
		Notes:     "test",
		CreatedAt: "2024-01-15T10:00:00Z",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if !containsKey(data, "item_name") {
		t.Error("JSON should contain 'item_name' key")
	}
	if !containsKey(data, "unit_price") {
		t.Error("JSON should contain 'unit_price' key")
	}
	if !containsKey(data, "unit_cost") {
		t.Error("JSON should contain 'unit_cost' key")
	}

	var decoded models.Sale
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.ItemName != original.ItemName {
		t.Errorf("ItemName mismatch: got %s, want %s", decoded.ItemName, original.ItemName)
	}
	if decoded.UnitPrice != original.UnitPrice {
		t.Errorf("UnitPrice mismatch: got %f, want %f", decoded.UnitPrice, original.UnitPrice)
	}
}

func TestUserJSON(t *testing.T) {
	original := models.User{
		ID:        "user-1",
		Username:  "testuser",
		Email:     "test@test.com",
		Password:  "hashedpassword",
		Role:      "vendor",
		ShopID:    "shop-1",
		Phone:     "123456789",
		Address:   "123 Test St",
		Bio:       "Test bio",
		AvatarURL: "/uploads/avatar.jpg",
		CreatedAt: "2024-01-15T10:00:00Z",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if !containsKey(data, "avatar_url") {
		t.Error("JSON should contain 'avatar_url' key")
	}
	if !containsKey(data, "shop_id") {
		t.Error("JSON should contain 'shop_id' key")
	}

	var decoded models.User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.Username != original.Username {
		t.Errorf("Username mismatch: got %s, want %s", decoded.Username, original.Username)
	}
	if decoded.AvatarURL != original.AvatarURL {
		t.Errorf("AvatarURL mismatch: got %s, want %s", decoded.AvatarURL, original.AvatarURL)
	}
}

func TestVendorJSON(t *testing.T) {
	original := models.Vendor{
		ID:        "vendor-1",
		Name:      "Test Vendor",
		Email:     "vendor@test.com",
		Role:      "vendor",
		ShopID:    "shop-1",
		CreatedAt: "2024-01-15T10:00:00Z",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if !containsKey(data, "shop_id") {
		t.Error("JSON should contain 'shop_id' key")
	}

	var decoded models.Vendor
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, original.Name)
	}
	if decoded.Email != original.Email {
		t.Errorf("Email mismatch: got %s, want %s", decoded.Email, original.Email)
	}
}

func TestPnLSummaryJSON(t *testing.T) {
	original := models.PnLSummary{
		TotalRevenue:  1000.0,
		TotalCOGS:     400.0,
		TotalExpenses: 200.0,
		GrossProfit:   600.0,
		NetProfit:     400.0,
		From:          "2024-01-01",
		To:            "2024-12-31",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if !containsKey(data, "total_revenue") {
		t.Error("JSON should contain 'total_revenue' key")
	}
	if !containsKey(data, "total_cogs") {
		t.Error("JSON should contain 'total_cogs' key")
	}
	if !containsKey(data, "gross_profit") {
		t.Error("JSON should contain 'gross_profit' key")
	}
	if !containsKey(data, "net_profit") {
		t.Error("JSON should contain 'net_profit' key")
	}

	var decoded models.PnLSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.NetProfit != original.NetProfit {
		t.Errorf("NetProfit mismatch: got %f, want %f", decoded.NetProfit, original.NetProfit)
	}
}

func TestUserJSONEmptyValues(t *testing.T) {
	original := models.User{}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var decoded models.User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.ID != "" {
		t.Errorf("expected empty ID, got '%s'", decoded.ID)
	}
	if decoded.Username != "" {
		t.Errorf("expected empty Username, got '%s'", decoded.Username)
	}
}

func TestPnLSummaryNetProfitCalculation(t *testing.T) {
	p := models.PnLSummary{
		TotalRevenue:  1000.0,
		TotalCOGS:     400.0,
		TotalExpenses: 200.0,
	}
	p.GrossProfit = p.TotalRevenue - p.TotalCOGS
	p.NetProfit = p.GrossProfit - p.TotalExpenses

	if p.GrossProfit != 600.0 {
		t.Errorf("expected gross profit 600.0, got %f", p.GrossProfit)
	}
	if p.NetProfit != 400.0 {
		t.Errorf("expected net profit 400.0, got %f", p.NetProfit)
	}
}

func containsKey(data []byte, key string) bool {
	return len(data) > 0 && stringContains(string(data), "\""+key+"\"")
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
