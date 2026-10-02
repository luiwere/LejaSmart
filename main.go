package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"LejaSmart/db"
	"LejaSmart/handlers"
)

func main() {
	// Load environment variables from .env if present (local dev)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system environment variables (normal in production/Render)")
	}

	// connect to Database
	db.Init()

	// profile pictures are written to disk, so make sure the directory exists
	if err := os.MkdirAll(handlers.UploadDir(), 0o755); err != nil {
		log.Println("Could not create upload directory:", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(handlers.UploadDir()))))

	http.HandleFunc("/register", handlers.RegisterPage)
	http.HandleFunc("/", handlers.DashboardPage)
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/vendor", handlers.VendorDashboard)
	http.HandleFunc("/vendor/expenses", handlers.VendorExpenses)
	http.HandleFunc("/vendor/inventory", handlers.VendorInventory)
	http.HandleFunc("/accountant", handlers.Accountantdashboard)
	http.HandleFunc("/logout", handlers.Logout)

	http.HandleFunc("/profile", handlers.ProfilePage)
	http.HandleFunc("/profile/data", handlers.ProfileData)
	http.HandleFunc("/profile/details", handlers.UpdateProfileDetails)
	http.HandleFunc("/profile/password", handlers.UpdateProfilePassword)
	http.HandleFunc("/profile/avatar", handlers.ProfileAvatar)

	http.HandleFunc("/me", handlers.Me)

	http.HandleFunc("/owner", handlers.OwnerDashboard)

	http.HandleFunc("/expenses", handlers.Expenses)
	http.HandleFunc("/expenses/", handlers.Expenses)
	http.HandleFunc("/inventory", handlers.Inventory)
	http.HandleFunc("/pnl", handlers.ProfitAndLoss)
	http.HandleFunc("/pnl/", handlers.ProfitAndLoss)
	http.HandleFunc("/sales", handlers.Sales)
	http.HandleFunc("/sales/", handlers.Sales)
	http.HandleFunc("/vendors", handlers.Vendors)

  // Get port from environment variable — required by Render
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // fallback for local development
    }

    fmt.Println("Server running on port", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
