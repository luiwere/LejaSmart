package main

import (
	"fmt"
	"log"
	"net/http"

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

	http.HandleFunc("/expenses", handlers.Expenses)
	http.HandleFunc("/inventory",handlers.Inventory)
	http.HandleFunc("/pnl", handlers.ProfitAndLoss)
	http.HandleFunc("/vendors", handlers.Vendors)

	fmt.Println("Sever running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080",nil))
}
