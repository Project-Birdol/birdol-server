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
)

var tokenData map[string]tokenInf

const (
	tokenIDsize        = 32
	tokenExpireSeconds = 604800 //604800[s]=1[week], int64

	deleteExpiredTokenIntervalSeconds = 86400 //86400[s]=1[day]
)

//tokenInf is a information of token linked by "userID"
type tokenInf struct {
	TokenID string    `json:"token_id"` //token id
	Updated time.Time `json:"updated"`  //last updated time
}

//SetToken creates(or update) and save new token, returns access token as string
func SetToken(sqldb *gorm.DB, userID uint) (string, error) {

	//create rondom token id
	token, err := generateRandomString(tokenIDsize)
	if err != nil {
		log.Println("failed to generate rondom string:", err)
		return "", err
	}

	//databaseに保存。SoftDeletedなコラムを見つけるためにUnscopedが必要
	if err := sqldb.Unscoped().Model(&model.AccessToken{}).Where("user_id = ?", userID).Updates(map[string]interface{}{"token": token, "token_updated": time.Now(), "deleted_at": gorm.Expr("NULL")}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := sqldb.Create(&model.AccessToken{UserID: userID, Token: token, TokenUpdated: time.Now()}).Error; err != nil {
				return "", err
			}
			return token, nil
		}
		return "", err
	}
	return token, nil
}

//DeleteToken delete stored access token
func DeleteToken(sqldb *gorm.DB, userID uint) error {

	//dbからTokenが保存されているか否か
	var c int64
	sqldb.Model(&model.AccessToken{}).Where("user_id=?", userID).Count(&c)
	if c == 0 {
		return gorm.ErrRecordNotFound
	}

	//delete
	if err := sqldb.Where("user_id=?", userID).Delete(&model.AccessToken{}).Error; err != nil {
		return err
	}

	return nil
}

//CheckToken checks if token is already stored in database: return error if not stored or already expired or mis-match token
func CheckToken(sqldb *gorm.DB, userID uint, token string) error {
	//dbからTokenが保存されているか否か
	var storedToken model.AccessToken
	if err := sqldb.Where("user_id = ?", userID).Select("user_id", "token", "token_updated").Take(&storedToken).Error; err != nil {
		return err
	}

	//tokenが一致しているかのチャック
	// TODO: Confirm DeviceID
	if storedToken.Token != token {
		//log.Println(errors.New(fmt.Sprintf("mis-match token: userID(%d), token(%s), storedtoken(%s)", userID, token, storedToken.Token)))
		return errors.New("token mismatch!")
	}

	//tokenの有効期限が切れていないかのチェック
	if err := checkTokenExpire(storedToken.TokenUpdated); err != nil {
		return err
	}

	//認証okの時にtokenの有効期限を伸ばす場合は、TokenUpdatedを現在時刻に変更する。
	if err := sqldb.Model(&model.AccessToken{}).Where("user_id=?", userID).Update("token_updated", time.Now()).Error; err != nil {
		return err
	}

	log.Printf("right match token: userID(%d), tokenID(%s)\n", userID, token)
	return nil

}

//checkTokenExpire checks if stored token is expired: if expired, return error
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

//StartDeleteExpiredTokens delete tokens if they are expired
func StartDeleteExpiredTokens() {
	go func() {
		for {
			time.Sleep(time.Second * deleteExpiredTokenIntervalSeconds)
			deleteAllExpiredtokens()
		}
	}()
}

//deleteAllExpiredtokens delete all expired tokens in database
func deleteAllExpiredtokens() {
	sqldb := database.SqlConnect()
	db, _ := sqldb.DB()
	defer db.Close()

	t := time.Now().Add(-1 * tokenExpireSeconds * time.Second)
	if err := sqldb.Where("token_updated < ?", t).Delete(&model.AccessToken{}); err != nil {
		log.Println(err)
	}

	log.Println("Delete all expired access tokens...")
}
