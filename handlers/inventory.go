package handlers

import (
	"LejaSmart/db"
	"encoding/json"
	"net/http"
)

func Inventory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	// GET - Fetch Invntory from a vendor
	case http.MethodGet:
		vendorID := r.URL.Query().Get("vendorID")
		items, err := db.GetInventory(getSessionRole(r), getSessionShopID(r), vendorID)
		if err != nil {
			http.Error(w, "Could not fetch inventory", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)

		// POST - add/update a supply item
	case http.MethodPost:
		var item struct {
			VendorID     string  `json:"vendor_id"`
			Name         string  `json:"name"`
			SupplierName string  `json:"supplier_name"`
			Status       string  `json:"status"`
			ReorderLevel float64 `json:"reorder_level"`
			ExpiryDate   string  `json:"expiry_date"`
			RestockedAt  string  `json:"restocked_at"`
			Quantity     float64 `json:"quantity"`
			Unit         string  `json:"unit"`
		}
		json.NewDecoder(r.Body).Decode(&item)
		err := db.SaveInventoryItem(
			getSessionRole(r),
			getSessionShopID(r),
			item.VendorID,
			item.Name,
			item.SupplierName,
			item.Status,
			item.ExpiryDate,
			item.RestockedAt,
			item.ReorderLevel,
			item.Quantity,
			item.Unit,
		)
		if err != nil {
			http.Error(w, "Could not add item to Inventory", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
