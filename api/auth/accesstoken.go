package auth

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"time"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	tokenIDsize        = 32
	tokenExpireSeconds = 604800 // 604800[s]=1[week], int64

	deleteExpiredTokenIntervalSeconds = 86400 // 86400[s]=1[day]
)

// SetToken creates(or update) and save new token, returns access token as string
func SetToken(userID uint, device_id string) (string, error) {

	// create rondom token id
	token, err := generateRandomString(tokenIDsize)
	if err != nil {
		log.Println("failed to generate rondom string:", err)
		return "", err
	}
	
	new_token := model.AccessToken {
		UserID: userID,
		DeviceID: device_id,
		Token: token,
		TokenUpdated: time.Now(),
	}

	// Use "ON DUPLICATE KEY UPDATE"
	if err := database.Sqldb.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"token": token, "token_updated": time.Now()}),
	}).Create(&new_token).Error; err != nil {
		return "", err
	}

	return token, nil
}

// DeleteToken delete stored access token
func DeleteToken(userID uint) error {

	// dbからTokenが保存されているか否か
	var c int64
	database.Sqldb.Model(&model.AccessToken{}).Where("user_id=?", userID).Count(&c)
	if c == 0 {
		return gorm.ErrRecordNotFound
	}

	// delete
	if err := database.Sqldb.Where("user_id=?", userID).Delete(&model.AccessToken{}).Error; err != nil {
		return err
	}

	return nil
}

// CheckToken checks if token is already stored in database: return error if not stored or already expired or mis-match token
func CheckToken(userID uint, device_id string, token string) error {
	// dbからTokenが保存されているか否か
	var storedTokens []model.AccessToken
	database.Sqldb.Where("user_id = ?", userID).Find(&storedTokens)
	if len(storedTokens) == 0 {
		return errors.New("tokens not found")
	}
	for i := 0; i < len(storedTokens); i++ {
		if storedTokens[i].Token != token {
			if i == len(storedTokens) - 1 { return errors.New("invalid token") }
			continue
		}
		if storedTokens[i].DeviceID != device_id {
			return errors.New("invalid device")
		}
		if err := checkTokenExpire(storedTokens[i].TokenUpdated); err != nil {
			return err;
		}
		break
	}

	// 認証okの時にtokenの有効期限を伸ばす場合は、TokenUpdatedを現在時刻に変更する。
	if err := database.Sqldb.Model(&model.AccessToken{}).Where("user_id = ? AND device_id = ?", userID, device_id).Update("token_updated", time.Now()).Error; err != nil {
		return err
	}

	log.Printf("right match token: user_id(%d), device_id(%s), token(%s)\n", userID, device_id, token)
	return nil
}

// checkTokenExpire checks if stored token is expired: if expired, return error
func checkTokenExpire(updated time.Time) error {
	if time.Since(updated).Seconds() > tokenExpireSeconds {
		return errors.New("token has expired")
	}
	return nil
}

// generateRandomString generate random string as access token
func generateRandomString(size int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	r := make([]byte, size)
	for i := 0; i < size; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		r[i] = letters[num.Int64()]
	}

	return string(r), nil
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
