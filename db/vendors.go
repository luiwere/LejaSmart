package db

import (
	"github.com/google/uuid"
	"Digiledger/models"
)

func GetAllVendors(role string) ([]models.Vendor, error) {
	conn := DBForRole(role)
	rows, err := conn.Query(`SELECT id,name,email,role,created_at FROM vendors`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		var v models.Vendor
		rows.Scan(&v.ID, &v.Name, &v.Email, &v.Role, &v.CreatedAt)
		vendors = append(vendors, v)
	}
	return vendors, nil
}

func CreateVendor(role, name, email, vendorRole string) error {
	id := uuid.New().String()
	conn := DBForRole(role)
	_, err := conn.Exec(
		`INSERT INTO vendors (id, name, email, role) VALUES (?,?,?,?)`,
		id, name, email, vendorRole,
	)
	return err
}

   