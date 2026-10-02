package db

import (
	"LejaSmart/models"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

var (
	ErrEmailTaken       = errors.New("email already taken")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
)

func CreateUser(username, email, password, role, shopName, shopCode string) (string, error) {
	id := uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	conn := DBForRole(role)
	// shopID is nil for owners, who are not attached to any shop.
	var shopID interface{}
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
	case "vendor":
		if shopName == "" {
			return "", errors.New("shop name is required")
		}
		shopID, generatedShopCode, err = createShop(conn, shopName)
		if err != nil {
			return "", err
		}
	case "owner":
		// Owners do not belong to a shop, so they carry no shop_id at all.
		shopID = nil
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

	// Create vendor record for vendor users so they appear in vendor lists
	if role == "vendor" {
		vendorID := uuid.New().String()
		_, err = conn.Exec(
			`INSERT INTO vendors (id, name, email, role, shop_id) VALUES (?, ?, ?, ?, ?)`,
			vendorID, username, email, "vendor", shopID,
		)
		if err != nil {
			return "", err
		}
	}

	return generatedShopCode, nil
}

const userColumns = `id, username, email, password, role, shop_id, phone, address, bio, avatar_url, created_at`

// scanUser scans the columns listed in userColumns into a User. Every optional
// column is nullable in the database, so it is read through sql.NullString.
func scanUser(scanner interface{ Scan(dest ...interface{}) error }) (models.User, error) {
	var u models.User
	var shopID, phone, address, bio, avatarURL, createdAt sql.NullString
	err := scanner.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &shopID,
		&phone, &address, &bio, &avatarURL, &createdAt)
	u.ShopID = shopID.String
	u.Phone = phone.String
	u.Address = address.String
	u.Bio = bio.String
	u.AvatarURL = avatarURL.String
	u.CreatedAt = createdAt.String
	return u, err
}

// userConns returns both databases in lookup order so profile reads and writes
// always reach the user regardless of which database they were created in.
func userConns(role string) []*DBConn {
	if role == "owner" {
		return []*DBConn{OwnerDB, DB}
	}
	return []*DBConn{DB, OwnerDB}
}

func GetUserByEmail(email string) (models.User, error) {
	var u models.User
	var err error
	for _, conn := range userConns("") {
		u, err = scanUser(conn.QueryRow(`SELECT `+userColumns+` FROM users WHERE email = ?`, email))
		if err == nil {
			return u, nil
		}
		if err != sql.ErrNoRows {
			return u, err
		}
	}
	return u, err
}

func GetUserByID(id string) (models.User, error) {
	var u models.User
	var err error
	for _, conn := range userConns("") {
		u, err = scanUser(conn.QueryRow(`SELECT `+userColumns+` FROM users WHERE id = ?`, id))
		if err == nil {
			return u, nil
		}
		if err != sql.ErrNoRows {
			return u, err
		}
	}
	return u, err
}

// ProfileUpdate carries the editable profile fields.
type ProfileUpdate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Bio      string `json:"bio"`
}

// UpdateProfile writes the editable profile fields for a user. The email is
// checked across both databases so it cannot collide with an existing account.
func UpdateProfile(id, role string, p ProfileUpdate) (models.User, error) {
	existing, err := GetUserByID(id)
	if err != nil {
		return existing, err
	}

	if p.Email != "" && p.Email != existing.Email {
		taken, err := GetUserByEmail(p.Email)
		if err == nil && taken.ID != id {
			return existing, ErrEmailTaken
		}
		if err != nil && err != sql.ErrNoRows {
			return existing, err
		}
	}

	email := p.Email
	if email == "" {
		email = existing.Email
	}

	username := p.Username
	if username == "" {
		username = existing.Username
	}

	for _, conn := range userConns(role) {
		res, err := conn.Exec(
			`UPDATE users SET username = ?, email = ?, phone = ?, address = ?, bio = ? WHERE id = ?`,
			username, email, p.Phone, p.Address, p.Bio, id,
		)
		if err != nil {
			return existing, err
		}
		if n, err := res.RowsAffected(); err == nil && n > 0 {
			// Keep the vendor listing in sync with the user it mirrors.
			if role == "vendor" {
				conn.Exec(`UPDATE vendors SET name = ?, email = ? WHERE email = ?`,
					username, email, existing.Email)
			}
			return GetUserByID(id)
		}
	}

	return existing, sql.ErrNoRows
}

// UpdateUserAvatar stores the avatar path for a user.
func UpdateUserAvatar(id, role, avatarURL string) error {
	for _, conn := range userConns(role) {
		res, err := conn.Exec(`UPDATE users SET avatar_url = ? WHERE id = ?`, avatarURL, id)
		if err != nil {
			return err
		}
		if n, err := res.RowsAffected(); err == nil && n > 0 {
			return nil
		}
	}
	return sql.ErrNoRows
}

// UpdateUserPassword re-hashes and stores a new password.
func UpdateUserPassword(id, role, newPassword string) error {
	if len(newPassword) < 6 {
		return ErrPasswordTooShort
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	for _, conn := range userConns(role) {
		res, err := conn.Exec(`UPDATE users SET password = ? WHERE id = ?`, string(hashed), id)
		if err != nil {
			return err
		}
		if n, err := res.RowsAffected(); err == nil && n > 0 {
			return nil
		}
	}
	return sql.ErrNoRows
}

func GetShopNameByID(shopID string) (string, error) {
	var name string
	err := DB.QueryRow(`SELECT name FROM shops WHERE id = ?`, shopID).Scan(&name)
	if err == sql.ErrNoRows {
		err = OwnerDB.QueryRow(`SELECT name FROM shops WHERE id = ?`, shopID).Scan(&name)
	}
	return name, err
}

func GetShopCodeByID(shopID string) (string, error) {
	var code string
	err := DB.QueryRow(`SELECT code FROM shops WHERE id = ?`, shopID).Scan(&code)
	if err == sql.ErrNoRows {
		err = OwnerDB.QueryRow(`SELECT code FROM shops WHERE id = ?`, shopID).Scan(&code)
	}
	return code, err
}

func getShopIDByCode(conn *DBConn, code string) (string, error) {
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

func createShop(conn *DBConn, name string) (string, string, error) {
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

func generateShopCode(conn *DBConn) (string, error) {
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
