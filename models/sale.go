package models

type Sale struct {
	ID        string  `json:"id"`
	VendorID  string  `json:"vendor_id"`
	ItemName  string  `json:"item_name"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	UnitCost  float64 `json:"unit_cost"`
	Date      string  `json:"date"`
	Notes     string  `json:"notes"`
	CreatedAt string  `json:"created_at"`
}
