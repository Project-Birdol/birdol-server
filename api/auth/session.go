package auth

import (
	"time"
	"crypto/sha256"
	"encoding/hex"

	"github.com/MISW/birdol-server/database/model"
	"gorm.io/gorm"
)

func CreateSession(db *gorm.DB, device_id string, access_token string) (string, error){
	var session model.Session
	if err := db.Where("device_id = ?", device_id).First(&session).Error; err != nil {
		
	}

	identifier := generateSessionIdentifier(device_id, access_token)

	return identifier, nil
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