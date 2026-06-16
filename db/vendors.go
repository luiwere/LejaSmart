package db

import (
	"github.com/google/uuid"
	"Digiledger/models"
)

func GetAllVendors() ([]models.Vendor, error) {
	rows, err := DB.Query(`SELECT id,name,email,role,created_at FROM vendors`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendors []models.Vendor 
	for rows.Next() {
		var v models.Vendor
		rows.Scan(&v.ID,&v.Name,&v.Email,&v.Role,&v.CreatedAt)
		vendors = append(vendors,v)
	}
	return vendors, nil
}

func CreateVendor(name, email, role string) error {
	id := uuid.New().String()
	_, err := DB.Exec (
		`INSERT INTO vendors (id, name, email, role) VALUES (?,?,?,?)`,
		id, name, email, role,
	)
	return err
}

   