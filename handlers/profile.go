package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"LejaSmart/db"
	"golang.org/x/crypto/bcrypt"
)

const maxAvatarBytes = 5 << 20 // 5 MB

var avatarExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// uploadDir is where profile pictures are written. Defaults to ./uploads so the
// files are kept out of the asset tree that the app ships.
func uploadDir() string {
	if dir := os.Getenv("UPLOAD_DIR"); dir != "" {
		return dir
	}
	return "uploads"
}

// UploadDir is the directory profile pictures are served from.
func UploadDir() string {
	return uploadDir()
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	if userID == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if _, err := db.GetUserByID(userID); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// profileResponse is the JSON shape shared by the read and write endpoints.
type profileResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	ShopID    string `json:"shop_id"`
	ShopName  string `json:"shop_name"`
	ShopCode  string `json:"shop_code"`
	CreatedAt string `json:"created_at"`
}

func buildProfileResponse(userID string) (profileResponse, error) {
	user, err := db.GetUserByID(userID)
	if err != nil {
		return profileResponse{}, err
	}

	shopName, shopCode := "", ""
	if user.ShopID != "" {
		if name, err := db.GetShopNameByID(user.ShopID); err == nil {
			shopName = name
		}
		if code, err := db.GetShopCodeByID(user.ShopID); err == nil {
			shopCode = code
		}
	}

	return profileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Phone:     user.Phone,
		Address:   user.Address,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		ShopID:    user.ShopID,
		ShopName:  shopName,
		ShopCode:  shopCode,
		CreatedAt: user.CreatedAt,
	}, nil
}

// ProfileData serves the signed-in user's profile as JSON.
func ProfileData(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	if userID == "" {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	profile, err := buildProfileResponse(userID)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfileDetails saves the editable profile fields.
func UpdateProfileDetails(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	if userID == "" {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	var p db.ProfileUpdate
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if _, err := db.UpdateProfile(userID, getSessionRole(r), p); err != nil {
		if err == db.ErrEmailTaken {
			http.Error(w, "Email already taken", http.StatusConflict)
			return
		}
		http.Error(w, "Could not save profile", http.StatusInternalServerError)
		return
	}

	profile, err := buildProfileResponse(userID)
	if err != nil {
		http.Error(w, "Could not load profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfilePassword changes the signed-in user's password.
func UpdateProfilePassword(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	if userID == "" {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	var p struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.CurrentPassword)); err != nil {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}
	if err := db.UpdateUserPassword(userID, getSessionRole(r), p.NewPassword); err != nil {
		if err == db.ErrPasswordTooShort {
			http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
			return
		}
		http.Error(w, "Could not update password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// ProfileAvatar handles upload and removal of the profile picture.
func ProfileAvatar(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	if userID == "" {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, maxAvatarBytes)
		if err := r.ParseMultipartForm(maxAvatarBytes); err != nil {
			http.Error(w, "Image must be 5MB or smaller", http.StatusBadRequest)
			return
		}
		file, header, err := r.FormFile("avatar")
		if err != nil {
			http.Error(w, "No image selected", http.StatusBadRequest)
			return
		}
		defer file.Close()

		ext, err := detectImageExtension(file, header)
		if err != nil {
			http.Error(w, "Only JPG, PNG, GIF or WEBP images are allowed", http.StatusBadRequest)
			return
		}

		dir := uploadDir()
		if err := os.MkdirAll(dir, 0o755); err != nil {
			http.Error(w, "Could not store image", http.StatusInternalServerError)
			return
		}

		name := fmt.Sprintf("%s%s", userID, ext)
		dst, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			http.Error(w, "Could not store image", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Could not store image", http.StatusInternalServerError)
			return
		}

		avatarURL := "/uploads/" + name
		if err := db.UpdateUserAvatar(userID, getSessionRole(r), avatarURL); err != nil {
			http.Error(w, "Could not save image", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"avatar_url": avatarURL})

	case http.MethodDelete:
		user, err := db.GetUserByID(userID)
		if err != nil {
			http.Error(w, "Not logged in", http.StatusUnauthorized)
			return
		}
		if user.AvatarURL != "" {
			os.Remove(filepath.Join(uploadDir(), filepath.Base(user.AvatarURL)))
		}
		if err := db.UpdateUserAvatar(userID, getSessionRole(r), ""); err != nil {
			http.Error(w, "Could not remove image", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"avatar_url": ""})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// detectImageExtension verifies the upload really is an image before trusting
// its filename, and returns the extension to store it under.
func detectImageExtension(file multipart.File, header *multipart.FileHeader) (string, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	if ext, ok := avatarExtensions[http.DetectContentType(buf[:n])]; ok {
		return ext, nil
	}
	// DetectContentType has no WEBP entry in some Go builds, so fall back to
	// the extension the browser reported.
	if ext := strings.ToLower(filepath.Ext(header.Filename)); ext == ".webp" {
		return ".webp", nil
	}
	return "", fmt.Errorf("unsupported image type")
}
