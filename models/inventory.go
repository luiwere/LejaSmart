package models

type InventoryItem struct {
    ID        string  `json:"id"`
    VendorID  string  `json:"vendor_id"`
    Name      string  `json:"name"`
    Quantity  float64 `json:"quantity"`
    Unit      string  `json:"unit"`
    UpdatedAt string  `json:"updated_at"`
}
