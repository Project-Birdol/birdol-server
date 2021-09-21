package middlewares

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"

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
		access_token := ctx.GetHeader("Authorization")
		device_id := ctx.GetHeader("device_id")
		signature := ctx.GetHeader("X-Birdol-Signature")
		timestamp := ctx.GetHeader("X-Birdol-TimeStamp")

		// confirm accesstoken
		var recv_token model.AccessToken
		if err := database.Sqldb.Where("access_token = ?", access_token).First(&recv_token).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H {
				"result": "failed",
				"error":  "Invalid token",
			})
			ctx.Abort()
		}

		// confirm device id
		if device_id != recv_token.DeviceID {
			ctx.JSON(http.StatusUnauthorized, gin.H {
				"result": "failed",
				"error":  "Invalid DeviceID",
			})
			ctx.Abort()
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
		}

		request_byte, _ := io.ReadAll(ctx.Request.Body)
		request_body := string(request_byte)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(request_body)))

		signature_base := "v2:" + timestamp + ":" + request_body
		hashed_base := sha256.Sum256([]byte(signature_base))

		verify_err := rsa.VerifyPKCS1v15(&public_key, crypto.SHA256, hashed_base[:], []byte(signature))
		if verify_err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H {
				"result": "failed",
				"error":  "Invalid signature",
			})
			ctx.Abort()
		}

		ctx.Set("AccessToken", recv_token)

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
