package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"Digiledger/db"
)

func Vendors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		vendors, err := db.GetAllVendors()
		if err != nil {
			http.Error(w, "Could not Fetch vendors", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vendors)

	case http.MethodPost:
		var v struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Role  string `json:"role"`
		}
		json.NewDecoder(r.Body).Decode(&v)
		err := db.CreateVendor(v.Name, v.Email, v.Role)
		if err != nil {
			http.Error(w, "Could not create Vendor", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func VendorDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/vendor-dashboard.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func Accountantdashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/accountant-dashboard.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
