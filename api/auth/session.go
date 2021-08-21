package auth

import (
	"time"
	"crypto/sha256"
	"encoding/hex"

	"github.com/MISW/birdol-server/database/model"
	"gorm.io/gorm"
)

func CreateSession(db *gorm.DB, device_id string, access_token string, user_id uint) (string, error) {
	db.Model(&model.Session{}).Where("user_id = ?", user_id).Update("expired", true)
	identifier := generateSessionIdentifier(device_id, access_token)
	new_session := model.Session {
		SessionID: identifier,
		AccessToken: access_token,
		UserID: user_id,
		Expired: false,
	}
	if err := db.Model(&model.Session{}).Where("device_id = ?", device_id).Updates(new_session).Error; err != nil {
		if err := db.Create(new_session).Error; err != nil {
			return "", err
		}
	}
	return identifier, nil
}

func CheckSession(db *gorm.DB, session_id string, access_token string) bool {
	var session model.Session
	if result := db.Model(&session).Where("session_id = ?", session_id);  result.Error != nil {
		return false
	}
	if session.AccessToken != access_token { return false }
	if session.Expired { return false }
	return true
}

func generateSessionIdentifier(device_id string, access_token string) string {
	t := time.Now().String()
	base_string := access_token + device_id + t
	hash_string := hex.EncodeToString(getBinarySHA256(base_string))
	return hash_string
}

func getBinarySHA256(base_string string) []byte {
	bin := sha256.Sum256([]byte(base_string))
	return bin[:]
}
