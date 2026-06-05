package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"Digiledger/db"
)

func ProfitAndLoss(w http.ResponseWriter, r *http.Request) {

	// Extract vendeoID from URL:/pnl/[vendorID]
	vendorID := strings.TrimPrefix(r.URL.Path, "/pnl")

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	summary, err := db.GetPnL(vendorID, from, to)

	if err != nil {
		http.Error(w, "Could Not Calculate Profit&Loss", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
