package db

import (
    "database/sql"
    "github.com/google/uuid"
    "Digiledger/models"
)

func GetExpenses(role, vendorID string) ([]models.Expense, error) {
    conn := DBForRole(role)
    var rows *sql.Rows
    var err error

    if vendorID == "" {
        rows, err = conn.Query(`SELECT id, vendor_id, amount, date, category, supplier_name, notes, created_at FROM expenses ORDER BY date DESC`)
    } else {
        rows, err = conn.Query(`SELECT id, vendor_id, amount, date, category, supplier_name, notes, created_at FROM expenses WHERE vendor_id = ? ORDER BY date DESC`, vendorID)
    }

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var expenses []models.Expense
    for rows.Next() {
        var e models.Expense
        rows.Scan(&e.ID, &e.VendorID, &e.Amount, &e.Date, &e.Category, &e.SupplierName, &e.Notes, &e.CreatedAt)
        expenses = append(expenses, e)
    }
    return expenses, nil
}

func AddExpense(role, vendorID string, amount float64, date, category, supplierName, notes string) error {
    conn := DBForRole(role)
    id := uuid.New().String()
    _, err := conn.Exec(
        `INSERT INTO expenses (id, vendor_id, amount, date, category, supplier_name, notes) VALUES (?, ?, ?, ?, ?, ?, ?)`,
        id, vendorID, amount, date, category, supplierName, notes,
    )
    return err
}

func DeleteExpense(role, id string) error {
    conn := DBForRole(role)
    _, err := conn.Exec(`DELETE FROM expenses WHERE id = ?`, id)
    return err
}
