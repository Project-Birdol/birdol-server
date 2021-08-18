package auth

import "golang.org/x/crypto/bcrypt"

//HashString string(パスワードなど)をハッシュ化する。
func HashString(s *string) error {
	bytes := []byte(*s)

	// `hash and set password
	hashed, err := bcrypt.GenerateFromPassword(bytes, 10) //(password, cost(4~31) )
	*s = string(hashed)

	return err
}

//CompareHashedString ハッシュ済stringと未ハッシュstringを比較し、一致した場合nilを返す。
func CompareHashedString(hashedPass string, rawPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(rawPass))
}
