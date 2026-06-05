package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"Digiledger/db"
)

func Expenses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		vendorID := r.URL.Query().Get("vendorID")
		expenses, err := db.GetExpenses(vendorID)
		if err != nil {
			http.Error(w, "Could not Fetch Expenses", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expenses)

	case http.MethodPost:
		var e struct {
			VendorID string `json:"vendor_id"`
			Amount float64 `json:"amount"`
			Date string `json:"date"`
			Category string `json:"category"`
			SupplierName string `json:"supplier_name"`
			Notes string `json:"notes"`
		}
		json.NewDecoder(r.Body).Decode(&e)
		err := db.AddExpense(e.VendorID, e.Amount, e.Date, e.Category, e.SupplierName, e.Notes)
		if err != nil {
			http.Error(w, "Could not save Expense", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		id := strings.TrimPrefix(r.URL.Path, "/expenses/")
		err := db.DeleteExpense(id)
		if err != nil {
			http.Error(w, "Could not delete expense", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
