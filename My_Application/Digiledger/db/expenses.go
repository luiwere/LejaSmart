package db

import (
    "github.com/google/uuid"
)

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

func GetExpenses(vendorID string) ([]Expense, error) {
    var rows *sql.Rows
    var err error

    if vendorID == "" {
        // Accountant — get all expenses
        rows, err = DB.Query(`SELECT id, vendor_id, amount, date, category, supplier_name, notes, created_at FROM expenses ORDER BY date DESC`)
    } else {
        // Vendor — get only their expenses
        rows, err = DB.Query(`SELECT id, vendor_id, amount, date, category, supplier_name, notes, created_at FROM expenses WHERE vendor_id = ? ORDER BY date DESC`, vendorID)
    }

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var expenses []Expense
    for rows.Next() {
        var e Expense
        rows.Scan(&e.ID, &e.VendorID, &e.Amount, &e.Date, &e.Category, &e.SupplierName, &e.Notes, &e.CreatedAt)
        expenses = append(expenses, e)
    }
    return expenses, nil
}

func AddExpense(vendorID string, amount float64, date, category, supplierName, notes string) error {
    id := uuid.New().String()
    _, err := DB.Exec(
        `INSERT INTO expenses (id, vendor_id, amount, date, category, supplier_name, notes) VALUES (?, ?, ?, ?, ?, ?, ?)`,
        id, vendorID, amount, date, category, supplierName, notes,
    )
    return err
}

func DeleteExpense(id string) error {
    _, err := DB.Exec(`DELETE FROM expenses WHERE id = ?`, id)
    return err
}
