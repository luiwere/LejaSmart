package test

import (
	"database/sql"
	"testing"

	"LejaSmart/db"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUserVendor(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	shopCode, err := db.CreateUser("testuser", "vendor@test.com", "password123", "vendor", "Test Shop", "")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if shopCode == "" {
		t.Error("expected non-empty shop code for vendor")
	}

	user, err := db.GetUserByEmail("vendor@test.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", user.Username)
	}
	if user.Role != "vendor" {
		t.Errorf("expected role 'vendor', got '%s'", user.Role)
	}
	if user.ShopID == "" {
		t.Error("expected non-empty shop_id for vendor")
	}

	shopName, err := db.GetShopNameByID(user.ShopID)
	if err != nil {
		t.Fatalf("GetShopNameByID failed: %v", err)
	}
	if shopName != "Test Shop" {
		t.Errorf("expected shop name 'Test Shop', got '%s'", shopName)
	}
}

func TestCreateUserAccountant(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestShop(t, "shop-1", "Accountant Shop", "ACCTSHOP")

	shopCode, err := db.CreateUser("acctuser", "acct@test.com", "password123", "accountant", "", "ACCTSHOP")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if shopCode != "" {
		t.Error("expected empty shop code for accountant")
	}

	user, err := db.GetUserByEmail("acct@test.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if user.ShopID != "shop-1" {
		t.Errorf("expected shop_id 'shop-1', got '%s'", user.ShopID)
	}
}

func TestCreateUserOwner(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	shopCode, err := db.CreateUser("owneruser", "owner@test.com", "password123", "owner", "", "")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if shopCode != "" {
		t.Error("expected empty shop code for owner")
	}

	user, err := db.GetUserByEmail("owner@test.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if user.ShopID != "" {
		t.Error("expected empty shop_id for owner")
	}
}

func TestCreateUserInvalidRole(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.CreateUser("baduser", "bad@test.com", "password123", "superadmin", "", "")
	if err == nil {
		t.Error("expected error for invalid role")
	}
}

func TestCreateUserVendorNoShopName(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.CreateUser("vendor", "v@test.com", "password123", "vendor", "", "")
	if err == nil {
		t.Error("expected error when vendor has no shop name")
	}
}

func TestCreateUserAccountantNoShopCode(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.CreateUser("acct", "a@test.com", "password123", "accountant", "", "")
	if err == nil {
		t.Error("expected error when accountant has no shop code")
	}
}

func TestGetUserByEmail(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestUser(t, "user-1", "john", "john@test.com", "vendor", "shop-1")

	user, err := db.GetUserByEmail("john@test.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if user.ID != "user-1" {
		t.Errorf("expected ID 'user-1', got '%s'", user.ID)
	}
	if user.Username != "john" {
		t.Errorf("expected username 'john', got '%s'", user.Username)
	}
}

func TestGetUserByEmailNotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.GetUserByEmail("nonexistent@test.com")
	if err == nil {
		t.Error("expected error for non-existent email")
	}
}

func TestGetUserByID(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestUser(t, "user-2", "jane", "jane@test.com", "owner", "")

	user, err := db.GetUserByID("user-2")
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if user.ID != "user-2" {
		t.Errorf("expected ID 'user-2', got '%s'", user.ID)
	}
	if user.Email != "jane@test.com" {
		t.Errorf("expected email 'jane@test.com', got '%s'", user.Email)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.GetUserByID("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestUpdateProfile(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestUser(t, "user-3", "bob", "bob@test.com", "vendor", "shop-1")

	updated, err := db.UpdateProfile("user-3", "vendor", db.ProfileUpdate{
		Username: "bob_updated",
		Phone:    "123456789",
		Address:  "123 Test St",
		Bio:      "Software developer",
	})
	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}
	if updated.Username != "bob_updated" {
		t.Errorf("expected username 'bob_updated', got '%s'", updated.Username)
	}
	if updated.Phone != "123456789" {
		t.Errorf("expected phone '123456789', got '%s'", updated.Phone)
	}
}

func TestUpdateProfileEmailTaken(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestUser(t, "user-a", "alice", "alice@test.com", "vendor", "shop-1")
	createTestUser(t, "user-b", "bob", "bob@test.com", "vendor", "shop-1")

	_, err := db.UpdateProfile("user-a", "vendor", db.ProfileUpdate{
		Email: "bob@test.com",
	})
	if err != db.ErrEmailTaken {
		t.Errorf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUpdateUserPassword(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestUser(t, "user-4", "charlie", "charlie@test.com", "vendor", "shop-1")

	err := db.UpdateUserPassword("user-4", "vendor", "newpassword456")
	if err != nil {
		t.Fatalf("UpdateUserPassword failed: %v", err)
	}

	user, err := db.GetUserByID("user-4")
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if user.Password == "hashedpassword" {
		t.Error("password hash should have changed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("newpassword456")); err != nil {
		t.Error("new password should match hash")
	}
}

func TestUpdateUserPasswordTooShort(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestUser(t, "user-5", "dave", "dave@test.com", "vendor", "shop-1")

	err := db.UpdateUserPassword("user-5", "vendor", "12345")
	if err != db.ErrPasswordTooShort {
		t.Errorf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestUpdateUserAvatar(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	createTestUser(t, "user-6", "eve", "eve@test.com", "vendor", "shop-1")

	err := db.UpdateUserAvatar("user-6", "vendor", "/uploads/avatar.jpg")
	if err != nil {
		t.Fatalf("UpdateUserAvatar failed: %v", err)
	}

	user, err := db.GetUserByID("user-6")
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if user.AvatarURL != "/uploads/avatar.jpg" {
		t.Errorf("expected avatar_url '/uploads/avatar.jpg', got '%s'", user.AvatarURL)
	}
}

func TestUpdateUserAvatarNotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	err := db.UpdateUserAvatar("nonexistent", "vendor", "/uploads/test.jpg")
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestCreateUserEmailUnique(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	_, err := db.CreateUser("user1", "same@test.com", "password123", "owner", "", "")
	if err != nil {
		t.Fatalf("first CreateUser failed: %v", err)
	}

	_, err = db.CreateUser("user2", "same@test.com", "password456", "owner", "", "")
	if err == nil {
		t.Error("expected error for duplicate email")
	}
}
