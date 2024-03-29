package auth

import (
	"log"
	"time"

	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/MISW/birdol-server/utils/random"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	tokenIDsize = 32
	refreshTokenSize = 16
	tokenExpireSeconds = 604800 // 604800[s] = 1[week], int64
	deleteExpiredTokenIntervalSeconds = 86400 // 86400[s] = 1[day]
)

// SetToken creates(or update) and save new token, returns access token as string
func SetToken(userID uint, device_id string, public_key string) (string, string, error) {

	// create rondom token id
	token, err := random.GenerateRandomString(tokenIDsize)
	if err != nil {
		log.Println("failed to generate rondom string:", err)
		return "", "", err
	}

	refresh_token, err := random.GenerateRandomString(tokenIDsize)
	if err != nil {
		log.Println("failed to generate rondom string:", err)
		return "", "", err
	}
	
	new_token := model.AccessToken {
		UserID: userID,
		DeviceID: device_id,
		Token: token,
		RefreshToken: refresh_token,
		TokenUpdated: time.Now(),
		PublicKey: public_key,
	}

	// Use "ON DUPLICATE KEY UPDATE"
	if err := database.Sqldb.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"token": token, "refresh_token": refresh_token, "token_updated": time.Now(), "public_key": public_key}),
	}).Create(&new_token).Error; err != nil {
		return "", "", err
	}

	return token, refresh_token, nil
}

// DeleteToken delete stored access token
func DeleteToken(token string, device_id string) error {

	// dbからTokenが保存されているか否か
	var c int64
	database.Sqldb.Model(&model.AccessToken{}).Where("token = ? AND device_id = ?", token, device_id).Count(&c)
	if c == 0 {
		return gorm.ErrRecordNotFound
	}

	// delete
	if err := database.Sqldb.Where("token = ? AND device_id = ?", token, device_id).Delete(&model.AccessToken{}).Error; err != nil {
		return err
	}

	return nil
}

// StartDeleteExpiredTokens delete tokens if they are expired
func StartDeleteExpiredTokens() {
	go func() {
		for {
			time.Sleep(time.Second * deleteExpiredTokenIntervalSeconds)
			deleteAllExpiredtokens()
		}
	}()
}

// deleteAllExpiredtokens delete all expired tokens in database
func deleteAllExpiredtokens() {

	t := time.Now().Add(-1 * tokenExpireSeconds * time.Second)
	if err := database.Sqldb.Where("token_updated < ?", t).Delete(&model.AccessToken{}); err != nil {
		log.Println(err)
	}
	log.Println("Delete all expired access tokens...")
}
