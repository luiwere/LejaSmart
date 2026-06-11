package models

type Expense struct {
    ID           string  `json:"id"`
    VendorID     string  `json:"vendor_id"`
    Amount       float64 `json:"amount"`
    Date         string  `json:"date"`
    Category     string  `json:"category"`
    SupplierName string  `json:"supplier_name"`
    Notes        string  `json:"notes"`
    CreatedAt    string  `json:"created_at"`
}
