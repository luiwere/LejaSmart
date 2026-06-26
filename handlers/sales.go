package handlers

import (
	"Digiledger/db"
	"encoding/json"
	"net/http"
	"strings"
)

func Sales(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		vendorID := r.URL.Query().Get("vendorID")
		sales, err := db.GetSales(getSessionRole(r), getSessionShopID(r), vendorID)
		if err != nil {
			http.Error(w, "Could not fetch sales", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sales)

	case http.MethodPost:
		var s struct {
			VendorID  string  `json:"vendor_id"`
			ItemName  string  `json:"item_name"`
			Quantity  float64 `json:"quantity"`
			UnitPrice float64 `json:"unit_price"`
			UnitCost  float64 `json:"unit_cost"`
			Date      string  `json:"date"`
			Notes     string  `json:"notes"`
		}
		json.NewDecoder(r.Body).Decode(&s)
		err := db.AddSale(getSessionRole(r), getSessionShopID(r), s.VendorID, s.ItemName, s.Quantity, s.UnitPrice, s.UnitCost, s.Date, s.Notes)
		if err != nil {
			http.Error(w, "Could not save sale", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		id := strings.TrimPrefix(r.URL.Path, "/sales/")
		err := db.DeleteSale(getSessionRole(r), getSessionShopID(r), id)
		if err != nil {
			http.Error(w, "Could not delete sale", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
