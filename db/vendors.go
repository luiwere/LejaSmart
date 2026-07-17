package db

import (
	"github.com/google/uuid"
	"LejaSmart/models"
)

func GetAllVendors(role, shopID string) ([]models.Vendor, error) {
	conn := DBForRole(role)
	rows, err := conn.Query(`SELECT id,name,email,role,shop_id,created_at FROM vendors WHERE shop_id = ?`, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		var v models.Vendor
		rows.Scan(&v.ID, &v.Name, &v.Email, &v.Role, &v.ShopID, &v.CreatedAt)
		vendors = append(vendors, v)
	}
	return vendors, nil
}

func CreateVendor(role, shopID, name, email, vendorRole string) error {
	id := uuid.New().String()
	conn := DBForRole(role)
	_, err := conn.Exec(
		`INSERT INTO vendors (id, name, email, role, shop_id) VALUES (?,?,?,?,?)`,
		id, name, email, vendorRole, shopID,
	)
	return err
}

   