package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/Project-Birdol/birdol-server/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type SessionManager struct {
	DB *gorm.DB
}

func (sm *SessionManager) CreateSession(device_id string, access_token string, user_id uint) (string, error) {
	sm.DB.Model(&model.Session{}).Where("user_id = ?", user_id).Update("expired", true)
	identifier := sm.generateSessionIdentifier(device_id, access_token)

	new_session := model.Session{
		SessionID:   identifier,
		AccessToken: access_token,
		UserID:      user_id,
		Expired:     false,
	}

	// Use "ON DUPLICATE KEY UPDATE"
	if err := sm.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "access_token"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"session_id": identifier, "expired": false}),
	}).Create(&new_session).Error; err != nil {
		return "", err
	}

	return identifier, nil
}

func (sm *SessionManager) CheckSession(session_id string, access_token string) bool {
	var session model.Session
	if result := sm.DB.Model(&session).Where("session_id = ?", session_id); result.Error != nil {
		return false
	}
	if session.AccessToken != access_token {
		return false
	}
	if session.Expired {
		return false
	}
	return true
}

func (sm *SessionManager) generateSessionIdentifier(device_id string, access_token string) string {
	t := time.Now().String()
	base_string := access_token + device_id + t
	hash_string := hex.EncodeToString(sm.getBinarySHA256(base_string))
	return hash_string
}

func (sm *SessionManager) getBinarySHA256(base_string string) []byte {
	bin := sha256.Sum256([]byte(base_string))
	return bin[:]
}
