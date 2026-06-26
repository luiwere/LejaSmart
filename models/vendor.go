package models

type Vendor struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    ShopID    string `json:"shop_id"`
    CreatedAt string `json:"created_at"`
}
