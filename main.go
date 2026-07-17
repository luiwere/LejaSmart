package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Digiledger/db"
	"Digiledger/handlers"
)

func main() {
	// connect to Database
	db.Init()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/register", handlers.RegisterPage)
	http.HandleFunc("/", handlers.LoginPage)
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/vendor", handlers.VendorDashboard)
	http.HandleFunc("/accountant", handlers.Accountantdashboard)
	http.HandleFunc("/logout", handlers.Logout)

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
