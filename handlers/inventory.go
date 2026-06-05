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
		items, err := db.GetInventory(VendorID)
		if err != nil {
			http.Error("Could not Fetch Inventoty", http.StatusSeerverInternalError)
			return
		}
		w.Header().Set("cotent-Type", "application/json")
		json.NewEncoder(w).encode(items)


		// POST - add/update a supply item
	case http.MethodPost:
		var item struct {
			VendorID string `json:"vendor_id"`
			Name string `json:"name"`
			Quantity float64 `json:"quantity"`
			Unit string `json:"unit"`
		}
		json.New.Decoder(r.Body).Decode(&item)
		err := db.SaveInventoryItem(item.VendorID, item.Name, item.Quantity, item.Unit)
		if err != nil {
		http.Error("Could not add item to Inventoyr", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
}


