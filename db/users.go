package db

import (
    "Digiledger/models"
    "github.com/google/uuid"
)

func CreateUser(username, email, password, role string) error {
    id := uuid.New().String()
    _, err := DB.Exec(
        `INSERT INTO users (id, username, email, password, role) VALUES (?, ?, ?, ?, ?)`,
        id, username, email, password, role,
    )
    return err
}

func GetUserByEmail(email string) (models.User, error) {
    var u models.User
    err := DB.QueryRow(
        `SELECT id, username, email, password, role, created_at FROM users WHERE email = ?`, email,
    ).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt)
    return u, err
}
