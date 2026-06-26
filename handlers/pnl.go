package handlers

import (
	"Digiledger/db"
	"encoding/json"
	"net/http"
	"strings"
)

func ProfitAndLoss(w http.ResponseWriter, r *http.Request) {

	// Extract vendorID from URL:/pnl/[vendorID]
	vendorID := strings.TrimPrefix(r.URL.Path, "/pnl")
	if strings.HasPrefix(vendorID, "/") {
		vendorID = strings.TrimPrefix(vendorID, "/")
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	summary, err := db.GetPnL(getSessionRole(r), vendorID, from, to)

	if err != nil {
		http.Error(w, "Could Not Calculate Profit&Loss", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
