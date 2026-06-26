package db

import (
    "database/sql"
    "errors"
    "strings"
    "Digiledger/models"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

func CreateUser(username, email, password, role, shopName, shopCode string) (string, error) {
    id := uuid.New().String()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }

    conn := DBForRole(role)
    shopID := ""
    generatedShopCode := ""

    switch role {
    case "accountant":
        if shopCode == "" {
            return "", errors.New("shop code is required for accountant")
        }
        shopID, err = getShopIDByCode(conn, shopCode)
        if err != nil {
            return "", err
        }
    case "vendor", "owner":
        if shopName == "" {
            return "", errors.New("shop name is required")
        }
        shopID, generatedShopCode, err = createShop(conn, shopName)
        if err != nil {
            return "", err
        }
    default:
        return "", errors.New("invalid role")
    }

    _, err = conn.Exec(
        `INSERT INTO users (id, username, email, password, role, shop_id) VALUES (?, ?, ?, ?, ?, ?)`,
        id, username, email, string(hashedPassword), role, shopID,
    )
    if err != nil {
        return "", err
    }
    return generatedShopCode, nil
}

func GetUserByEmail(email string) (models.User, error) {
    var u models.User
    conn := DBForEmail(email)
    err := conn.QueryRow(
        `SELECT id, username, email, password, role, shop_id, created_at FROM users WHERE email = ?`, email,
    ).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.ShopID, &u.CreatedAt)
    if err == sql.ErrNoRows && conn == DB {
        conn = OwnerDB
        err = conn.QueryRow(`SELECT id, username, email, password, role, shop_id, created_at FROM users WHERE email = ?`, email).
            Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.ShopID, &u.CreatedAt)
    }
    return u, err
}

func GetUserByID(id string) (models.User, error) {
    var u models.User
    err := DB.QueryRow(`SELECT id, username, email, password, role, shop_id, created_at FROM users WHERE id = ?`, id).
        Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.ShopID, &u.CreatedAt)
    if err == sql.ErrNoRows {
        err = OwnerDB.QueryRow(`SELECT id, username, email, password, role, shop_id, created_at FROM users WHERE id = ?`, id).
            Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.ShopID, &u.CreatedAt)
    }
    return u, err
}

func getShopIDByCode(conn *sql.DB, code string) (string, error) {
    var shopID string
    err := conn.QueryRow(`SELECT id FROM shops WHERE code = ?`, code).Scan(&shopID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", errors.New("shop code not found")
        }
        return "", err
    }
    return shopID, nil
}

func createShop(conn *sql.DB, name string) (string, string, error) {
    shopCode, err := generateShopCode(conn)
    if err != nil {
        return "", "", err
    }

    shopID := uuid.New().String()
    _, err = conn.Exec(
        `INSERT INTO shops (id, name, code) VALUES (?, ?, ?)`,
        shopID, name, shopCode,
    )
    if err != nil {
        return "", "", err
    }
    return shopID, shopCode, nil
}

func generateShopCode(conn *sql.DB) (string, error) {
    for i := 0; i < 10; i++ {
        code := strings.ToUpper(strings.ReplaceAll(uuid.New().String()[:8], "-", ""))
        if _, err := getShopIDByCode(conn, code); err != nil {
            if err.Error() == "shop code not found" {
                return code, nil
            }
            return "", err
        }
    }
    return "", errors.New("could not generate unique shop code")
}
