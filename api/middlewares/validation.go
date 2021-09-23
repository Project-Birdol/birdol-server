package middlewares

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"io"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/gin-gonic/gin"
)

// PublicKey XML structure
type rsaPublicKey struct {
	Modulus	string
	Exponent string
}

/*
	Main function
*/

func RequestValidation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// pick params from header
		authorization := ctx.GetHeader("Authorization")
		device_id := ctx.GetHeader("device_id")
		signature_str := ctx.GetHeader("X-Birdol-Signature")
		timestamp := ctx.GetHeader("X-Birdol-TimeStamp")

		// Verify Authorization Header
		reg := regexp.MustCompile(`Bearer (.+)$`)
		if !reg.MatchString(authorization) {
			ctx.JSON(http.StatusUnauthorized, gin.H {
				"result": "failed",
				"error": "invalid_token",
			})
			ctx.Abort()
			return
		}
		access_token := reg.FindStringSubmatch(authorization)[1]

		// confirm accesstoken
		var recv_token model.AccessToken
		if err := database.Sqldb.Where("token = ?", access_token).First(&recv_token).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H {
				"result": "failed",
				"error":  "invalid_token",
			})
			ctx.Abort()
			return
		}

		// confirm device id
		if device_id != recv_token.DeviceID {
			ctx.JSON(http.StatusUnauthorized, gin.H {
				"result": "failed",
				"error":  "invalid_DeviceID",
			})
			ctx.Abort()
			return
		}

		// Verify signature
		encoded_xml_key := recv_token.PublicKey
		public_key, err := mappingXML(encoded_xml_key)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"result": "failed",
				"error": "something went wrong",
			})
			ctx.Abort()
			return
		}

		body_byte, _ := io.ReadAll(ctx.Request.Body)
		ctx.Set("body_rawbyte", body_byte)
		request_body := string(body_byte)

		prefix := os.Getenv("API_VERSION")
		replacer := strings.NewReplacer("\r\n", "\n")
		request_body = replacer.Replace(request_body)
		signature_base := prefix + ":" + timestamp + ":" + request_body
		hashed_base := sha256.Sum256([]byte(signature_base))
		signature, _ := hex.DecodeString(signature_str)

		verify_err := rsa.VerifyPKCS1v15(&public_key, crypto.SHA256, hashed_base[:], signature)
		if verify_err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H {
				"result": "failed",
				"error":  "Invalid signature",
			})
			ctx.Abort()
			return
		}

		ctx.Set("access_token", recv_token)

		ctx.Next()
	}
}

/*
	private functions
*/

// decode base64 encoded string
func base64Decode(str string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// convert to bigInt value from base64 encoded string
func convertbigInt(str string) (*big.Int, error) {
	bigInt := &big.Int{}
	rawbyte, err := base64Decode(str)
	if err != nil {
		return nil, err
	}
	bigInt.SetBytes(rawbyte)
	return bigInt, nil
}

// map PublicKey XML to rsa.PublicKey struct
func mappingXML(str string) (rsa.PublicKey, error) {
	rawXML, _ := base64.StdEncoding.DecodeString(str)
	rsaPublicKey := rsaPublicKey{}
	if err := xml.Unmarshal([]byte(rawXML), &rsaPublicKey); err != nil {
		return rsa.PublicKey{}, err
	}
	modulus, err := convertbigInt(rsaPublicKey.Modulus)
	if err != nil { return rsa.PublicKey{}, err }
	exponent, err := convertbigInt(rsaPublicKey.Exponent)
	if err != nil { return rsa.PublicKey{}, err }
	key := rsa.PublicKey {
		N: modulus,
		E: int(exponent.Int64()),
	}
	return key, nil
}
