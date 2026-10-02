package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"LejaSmart/db"
	"LejaSmart/handlers"
)

func TestExpensesHandlerGET(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	db.AddExpense("vendor", "shop-1", "vendor-1", 100.0, "2024-01-15", "food", "Sup A", "")

	req := newAuthenticatedRequest(http.MethodGet, "/expenses", "vendor-1", "vendor", nil)
	w := httptest.NewRecorder()
	handlers.Expenses(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var expenses []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&expenses); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(expenses) != 1 {
		t.Errorf("expected 1 expense, got %d", len(expenses))
	}
}

func TestExpensesHandlerPOST(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	body := `{"vendor_id":"vendor-1","amount":250.50,"date":"2024-01-15","category":"food","supplier_name":"Supplier A","notes":"test"}`
	req := newAuthenticatedRequest(http.MethodPost, "/expenses", "vendor-1", "vendor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlers.Expenses(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	expenses, _ := db.GetExpenses("vendor", "shop-1", "")
	if len(expenses) != 1 {
		t.Errorf("expected 1 expense in DB, got %d", len(expenses))
	}
	if expenses[0].Amount != 250.50 {
		t.Errorf("expected amount 250.50, got %f", expenses[0].Amount)
	}
}

func TestExpensesHandlerDELETE(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	db.AddExpense("vendor", "shop-1", "vendor-1", 100.0, "2024-01-15", "food", "Sup A", "")

	expenses, _ := db.GetExpenses("vendor", "shop-1", "")
	if len(expenses) != 1 {
		t.Fatalf("expected 1 expense before delete, got %d", len(expenses))
	}

	req := newAuthenticatedRequest(http.MethodDelete, "/expenses/"+expenses[0].ID, "vendor-1", "vendor", nil)
	w := httptest.NewRecorder()
	handlers.Expenses(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	after, _ := db.GetExpenses("vendor", "shop-1", "")
	if len(after) != 0 {
		t.Errorf("expected 0 expenses after delete, got %d", len(after))
	}
}

func TestInventoryHandlerGET(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	db.SaveInventoryItem("vendor", "shop-1", "vendor-1", "Item A", "Sup A", "In Stock", "2024-12-31", "2024-01-01", 5.0, 100.0, "kg")

	req := newAuthenticatedRequest(http.MethodGet, "/inventory?vendorID=vendor-1", "vendor-1", "vendor", nil)
	w := httptest.NewRecorder()
	handlers.Inventory(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var items []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
}

func TestInventoryHandlerPOST(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	body := `{"vendor_id":"vendor-1","name":"New Item","supplier_name":"Supplier","status":"In Stock","reorder_level":5,"expiry_date":"2024-12-31","restocked_at":"2024-01-01","quantity":50.0,"unit":"kg"}`
	req := newAuthenticatedRequest(http.MethodPost, "/inventory", "vendor-1", "vendor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlers.Inventory(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	items, _ := db.GetInventory("vendor", "shop-1", "vendor-1")
	if len(items) != 1 {
		t.Errorf("expected 1 item in DB, got %d", len(items))
	}
	if items[0].Name != "New Item" {
		t.Errorf("expected name 'New Item', got '%s'", items[0].Name)
	}
}

func TestSalesHandlerGET(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	db.AddSale("vendor", "shop-1", "vendor-1", "Sugar", 10.0, 50.0, 30.0, "2024-01-15", "")

	req := newAuthenticatedRequest(http.MethodGet, "/sales", "vendor-1", "vendor", nil)
	w := httptest.NewRecorder()
	handlers.Sales(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var sales []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&sales); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(sales) != 1 {
		t.Errorf("expected 1 sale, got %d", len(sales))
	}
}

func TestSalesHandlerPOST(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	body := `{"vendor_id":"vendor-1","item_name":"New Item","quantity":5.0,"unit_price":100.0,"unit_cost":50.0,"date":"2024-01-15","notes":"test"}`
	req := newAuthenticatedRequest(http.MethodPost, "/sales", "vendor-1", "vendor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlers.Sales(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	sales, _ := db.GetSales("vendor", "shop-1", "")
	if len(sales) != 1 {
		t.Errorf("expected 1 sale in DB, got %d", len(sales))
	}
	if sales[0].ItemName != "New Item" {
		t.Errorf("expected item_name 'New Item', got '%s'", sales[0].ItemName)
	}
}

func TestSalesHandlerDELETE(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	db.AddSale("vendor", "shop-1", "vendor-1", "Sugar", 10.0, 50.0, 30.0, "2024-01-15", "")

	sales, _ := db.GetSales("vendor", "shop-1", "")
	if len(sales) != 1 {
		t.Fatalf("expected 1 sale before delete, got %d", len(sales))
	}

	req := newAuthenticatedRequest(http.MethodDelete, "/sales/"+sales[0].ID, "vendor-1", "vendor", nil)
	w := httptest.NewRecorder()
	handlers.Sales(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	after, _ := db.GetSales("vendor", "shop-1", "")
	if len(after) != 0 {
		t.Errorf("expected 0 sales after delete, got %d", len(after))
	}
}

func TestVendorsHandlerGET(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	db.CreateVendor("vendor", "shop-1", "Vendor A", "a@test.com", "vendor")
	db.CreateVendor("vendor", "shop-1", "Vendor B", "b@test.com", "vendor")

	req := newAuthenticatedRequest(http.MethodGet, "/vendors", "vendor-1", "vendor", nil)
	w := httptest.NewRecorder()
	handlers.Vendors(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var vendors []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&vendors); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(vendors) != 2 {
		t.Errorf("expected 2 vendors, got %d", len(vendors))
	}
}

func TestVendorsHandlerPOST(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	body := `{"name":"New Vendor","email":"new@test.com","role":"vendor"}`
	req := newAuthenticatedRequest(http.MethodPost, "/vendors", "vendor-1", "vendor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handlers.Vendors(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	vendors, _ := db.GetAllVendors("vendor", "shop-1")
	if len(vendors) != 1 {
		t.Errorf("expected 1 vendor in DB, got %d", len(vendors))
	}
}

func TestProfitAndLossHandler(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	db.AddExpense("vendor", "shop-1", "vendor-1", 800.0, "2024-01-15", "food", "Sup A", "")
	db.AddSale("vendor", "shop-1", "vendor-1", "Sugar", 10.0, 80.0, 40.0, "2024-01-15", "")

	req := newAuthenticatedRequest(http.MethodGet, "/pnl", "vendor-1", "vendor", nil)
	w := httptest.NewRecorder()
	handlers.ProfitAndLoss(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var summary map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&summary); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if summary["total_expenses"] == nil {
		t.Error("expected total_expenses in response")
	}
}

func TestMeHandler(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "user-1", "testuser", "user@test.com", "vendor", "shop-1")

	req := newAuthenticatedRequest(http.MethodGet, "/me", "user-1", "vendor", nil)
	w := httptest.NewRecorder()
	handlers.Me(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["username"] != "testuser" {
		t.Errorf("expected username 'testuser', got '%v'", resp["username"])
	}
	if resp["role"] != "vendor" {
		t.Errorf("expected role 'vendor', got '%v'", resp["role"])
	}
	if resp["shop_name"] != "Test Shop" {
		t.Errorf("expected shop_name 'Test Shop', got '%v'", resp["shop_name"])
	}
}

func TestMeHandlerNoAuth(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	handlers.Me(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestLogoutHandler(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()
	handlers.Logout(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, w.Code)
	}

	resp := w.Result()
	cookies := resp.Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "session_user" && c.Value == "" {
			found = true
			break
		}
	}
	if !found {
		t.Error("session_user cookie should be cleared in response")
	}
}

func TestProfitAndLossHandlerByVendor(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	createTestUser(t, "vendor-2", "v2", "v2@test.com", "vendor", "shop-1")

	db.AddExpense("vendor", "shop-1", "vendor-1", 300.0, "2024-01-15", "food", "Sup A", "")
	db.AddExpense("vendor", "shop-1", "vendor-2", 200.0, "2024-01-15", "food", "Sup B", "")
	db.AddSale("vendor", "shop-1", "vendor-1", "Item A", 5.0, 50.0, 25.0, "2024-01-15", "")
	db.AddSale("vendor", "shop-1", "vendor-2", "Item B", 3.0, 80.0, 40.0, "2024-01-15", "")

	req := newAuthenticatedRequest(http.MethodGet, "/pnl/vendor-1", "vendor-1", "vendor", nil)
	w := httptest.NewRecorder()
	handlers.ProfitAndLoss(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var summary map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&summary); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if summary["total_expenses"] != float64(300) {
		t.Errorf("expected total_expenses 300, got %v", summary["total_expenses"])
	}
	if summary["total_revenue"] != float64(250) {
		t.Errorf("expected total_revenue 250, got %v", summary["total_revenue"])
	}
}
