package test

import (
	"testing"

	"LejaSmart/db"
	"LejaSmart/models"
)

func TestGetPnLNoData(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")

	summary, err := db.GetPnL("vendor", "shop-1", "", "", "")
	if err != nil {
		t.Fatalf("GetPnL failed: %v", err)
	}
	if summary.TotalRevenue != 0 {
		t.Errorf("expected total revenue 0, got %f", summary.TotalRevenue)
	}
	if summary.TotalExpenses != 0 {
		t.Errorf("expected total expenses 0, got %f", summary.TotalExpenses)
	}
	if summary.GrossProfit != 0 {
		t.Errorf("expected gross profit 0, got %f", summary.GrossProfit)
	}
	if summary.NetProfit != 0 {
		t.Errorf("expected net profit 0, got %f", summary.NetProfit)
	}
}

func TestGetPnLWithData(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	db.AddExpense("vendor", "shop-1", "vendor-1", 500.0, "2024-01-10", "food", "Supplier A", "")
	db.AddExpense("vendor", "shop-1", "vendor-1", 300.0, "2024-01-15", "transport", "Supplier B", "")

	db.AddSale("vendor", "shop-1", "vendor-1", "Sugar", 10.0, 80.0, 40.0, "2024-01-12", "")
	db.AddSale("vendor", "shop-1", "vendor-1", "Salt", 20.0, 50.0, 20.0, "2024-01-20", "")

	summary, err := db.GetPnL("vendor", "shop-1", "", "", "")
	if err != nil {
		t.Fatalf("GetPnL failed: %v", err)
	}

	expectedExpenses := 800.0
	if summary.TotalExpenses != expectedExpenses {
		t.Errorf("expected total expenses %f, got %f", expectedExpenses, summary.TotalExpenses)
	}

	expectedRevenue := 10.0*80.0 + 20.0*50.0
	if summary.TotalRevenue != expectedRevenue {
		t.Errorf("expected total revenue %f, got %f", expectedRevenue, summary.TotalRevenue)
	}

	expectedCOGS := 10.0*40.0 + 20.0*20.0
	if summary.TotalCOGS != expectedCOGS {
		t.Errorf("expected total COGS %f, got %f", expectedCOGS, summary.TotalCOGS)
	}

	expectedGrossProfit := expectedRevenue - expectedCOGS
	if summary.GrossProfit != expectedGrossProfit {
		t.Errorf("expected gross profit %f, got %f", expectedGrossProfit, summary.GrossProfit)
	}

	expectedNetProfit := expectedGrossProfit - expectedExpenses
	if summary.NetProfit != expectedNetProfit {
		t.Errorf("expected net profit %f, got %f", expectedNetProfit, summary.NetProfit)
	}
}

func TestGetPnLDateFilter(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")

	db.AddExpense("vendor", "shop-1", "vendor-1", 500.0, "2024-01-10", "food", "Supplier A", "")
	db.AddExpense("vendor", "shop-1", "vendor-1", 300.0, "2024-06-15", "transport", "Supplier B", "")
	db.AddSale("vendor", "shop-1", "vendor-1", "Sugar", 10.0, 80.0, 40.0, "2024-01-12", "")
	db.AddSale("vendor", "shop-1", "vendor-1", "Salt", 20.0, 50.0, 20.0, "2024-06-20", "")

	summary, err := db.GetPnL("vendor", "shop-1", "", "2024-01-01", "2024-03-31")
	if err != nil {
		t.Fatalf("GetPnL failed: %v", err)
	}

	if summary.TotalExpenses != 500.0 {
		t.Errorf("expected total expenses 500, got %f", summary.TotalExpenses)
	}
	if summary.TotalRevenue != 800.0 {
		t.Errorf("expected total revenue 800, got %f", summary.TotalRevenue)
	}

	summary2, err := db.GetPnL("vendor", "shop-1", "", "2024-04-01", "2024-12-31")
	if err != nil {
		t.Fatalf("GetPnL failed: %v", err)
	}

	if summary2.TotalExpenses != 300.0 {
		t.Errorf("expected total expenses 300, got %f", summary2.TotalExpenses)
	}
	if summary2.TotalRevenue != 1000.0 {
		t.Errorf("expected total revenue 1000, got %f", summary2.TotalRevenue)
	}
}

func TestGetPnLByVendor(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Test Shop", "CODE1")
	createTestUser(t, "vendor-1", "v1", "v1@test.com", "vendor", "shop-1")
	createTestUser(t, "vendor-2", "v2", "v2@test.com", "vendor", "shop-1")

	db.AddExpense("vendor", "shop-1", "vendor-1", 100.0, "2024-01-10", "food", "Sup A", "")
	db.AddExpense("vendor", "shop-1", "vendor-2", 200.0, "2024-01-15", "food", "Sup B", "")
	db.AddSale("vendor", "shop-1", "vendor-1", "Item A", 5.0, 50.0, 25.0, "2024-01-12", "")
	db.AddSale("vendor", "shop-1", "vendor-2", "Item B", 3.0, 100.0, 50.0, "2024-01-16", "")

	summary, err := db.GetPnL("vendor", "shop-1", "vendor-1", "", "")
	if err != nil {
		t.Fatalf("GetPnL failed: %v", err)
	}
	if summary.TotalExpenses != 100.0 {
		t.Errorf("expected total expenses 100, got %f", summary.TotalExpenses)
	}
	if summary.TotalRevenue != 250.0 {
		t.Errorf("expected total revenue 250, got %f", summary.TotalRevenue)
	}

	summary2, err := db.GetPnL("vendor", "shop-1", "vendor-2", "", "")
	if err != nil {
		t.Fatalf("GetPnL failed: %v", err)
	}
	if summary2.TotalExpenses != 200.0 {
		t.Errorf("expected total expenses 200, got %f", summary2.TotalExpenses)
	}
	if summary2.TotalRevenue != 300.0 {
		t.Errorf("expected total revenue 300, got %f", summary2.TotalRevenue)
	}
}

func TestPnLSummaryModel(t *testing.T) {
	p := models.PnLSummary{
		TotalRevenue:  1000.0,
		TotalCOGS:     400.0,
		TotalExpenses: 200.0,
		GrossProfit:   600.0,
		NetProfit:     400.0,
		From:          "2024-01-01",
		To:            "2024-12-31",
	}
	if p.NetProfit != 400.0 {
		t.Errorf("expected net profit 400.0, got %f", p.NetProfit)
	}
}
