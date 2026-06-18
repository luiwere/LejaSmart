package handlers

import (
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
        email       := r.FormValue("email")
        password    := r.FormValue("password")

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
    } else {
        http.Redirect(w, r, "/vendor", http.StatusSeeOther)
    }
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
        email    := r.FormValue("email")
        password := r.FormValue("password")
        role     := r.FormValue("role")

        if username == "" || email == "" || password == "" {
            http.Error(w, "All fields are required", http.StatusBadRequest)
            return
        }

        err := db.CreateUser(username, email, password, role)
        if err != nil {
            if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
                http.Error(w, "Email already taken", http.StatusConflict)
            } else {
                http.Error(w, "Could not create account", http.StatusInternalServerError)
            }
            return
        }
    }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }

func Logout(w http.ResponseWriter, r *http.Request) {

    http.SetCookie(w, &http.Cookie{
        Name:     "session_user",
        Value:    "",
        Expires:  time.Now().Add(-1 * time.Hour),
        Path:     "/",
    })

    http.SetCookie(w, &http.Cookie{
        Name:     "session_role",
        Value:    "",
        Expires:  time.Now().Add(-1 * time.Hour),
        Path:     "/",
    })
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
