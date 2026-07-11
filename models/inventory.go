package models

type InventoryItem struct {
	ID           string  `json:"id"`
	VendorID     string  `json:"vendor_id"`
	Name         string  `json:"name"`
	SupplierName string  `json:"supplier_name"`
	Status       string  `json:"status"`
	ReorderLevel float64 `json:"reorder_level"`
	ExpiryDate   string  `json:"expiry_date"`
	RestockedAt  string  `json:"restocked_at"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	UpdatedAt    string  `json:"updated_at"`
}
