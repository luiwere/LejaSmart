package handlers

import (
	"html/template"
	"net/http"
	"Digiledger/db"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        http.ServeFile(w, r, "templates/register.html")
        return
    }

    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        email    := r.FormValue("email")
        password := r.FormValue("password")
        role     := r.FormValue("role")

        if username == "" || email == "" || password == "" {
            http.Error(w, "All fields are required", http.StatusBadRequest)
            return
        }

        err := db.CreateUser(username, email, password, role)
        if err != nil {
            http.Error(w, "Email already taken", http.StatusConflict)
            return
        }

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
