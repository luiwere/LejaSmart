package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"Digiledger/db"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Could not load page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		user, err := db.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "Invalid Email or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid Email or password", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_user",
			Value:    user.ID,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "session_role",
			Value:    user.Role,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})

		if user.Role == "accountant" {
			http.Redirect(w, r, "/accountant", http.StatusSeeOther)
			return
		} else if user.Role == "owner" {
			http.Redirect(w, r, "/owner", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/vendor", http.StatusSeeOther)
		return
	}
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Could not load page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		role := r.FormValue("role")
		shopName := r.FormValue("shop_name")
		shopCode := r.FormValue("shop_code")

		if username == "" || email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if role == "accountant" && shopCode == "" {
			http.Error(w, "Shop code is required for accountant", http.StatusBadRequest)
			return
		}

		if role == "vendor" && shopName == "" {
			http.Error(w, "Shop name is required for vendor", http.StatusBadRequest)
			return
		}

		generatedShopCode, err := db.CreateUser(username, email, password, role, shopName, shopCode)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				http.Error(w, "Email already taken", http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if generatedShopCode != "" {
			json.NewEncoder(w).Encode(map[string]string{"shop_code": generatedShopCode})
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
		return
	}
}

func Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_user")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	user, err := db.GetUserByID(cookie.Value)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	shopName := ""
	shopCode := ""
	if user.ShopID != "" {
		name, err := db.GetShopNameByID(user.ShopID)
		if err == nil {
			shopName = name
		}
		code, err := db.GetShopCodeByID(user.ShopID)
		if err == nil {
			shopCode = code
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":        user.ID,
		"role":      user.Role,
		"shop_id":   user.ShopID,
		"shop_name": shopName,
		"shop_code": shopCode,
	})
}

func getSessionUserID(r *http.Request) string {
	cookie, err := r.Cookie("session_user")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func getSessionRole(r *http.Request) string {
	roleCookie, err := r.Cookie("session_role")
	if err != nil {
		return ""
	}
	return roleCookie.Value
}

func getSessionShopID(r *http.Request) string {
	userID := getSessionUserID(r)
	if userID == "" {
		return ""
	}
	user, err := db.GetUserByID(userID)
	if err != nil {
		return ""
	}
	return user.ShopID
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_user",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "session_role",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
