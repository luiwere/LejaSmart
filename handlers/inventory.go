package handlers

import (
	"encoding/json"
	"net/http"
	"Digiledger/db"

)

func Inventory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

		// GET - Fetch Invntory from a vendor
	case http.MethodGet:
		vendorID := r.URL.Query().Get("vendorID")
		items, err := db.GetInventory(getSessionRole(r), vendorID)
		if err != nil {
			http.Error(w, "Could not Fetch Inventoty", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Cotent-Type", "application/json")
		json.NewEncoder(w).Encode(items)


		// POST - add/update a supply item
	case http.MethodPost:
		var item struct {
			VendorID string `json:"vendor_id"`
			Name string `json:"name"`
			Quantity float64 `json:"quantity"`
			Unit string `json:"unit"`
		}
		json.NewDecoder(r.Body).Decode(&item)
		err := db.SaveInventoryItem(getSessionRole(r), item.VendorID, item.Name, item.Quantity, item.Unit)
		if err != nil {
		http.Error(w , "Could not add item to Inventoyr", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
}


