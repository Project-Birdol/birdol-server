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
	"net/http/httputil"
	"log"
	"os"
	"regexp"
	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/MISW/birdol-server/utils/response"
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
		// RequestDump
		requestDump, err := httputil.DumpRequest(ctx.Request, true)
		if err != nil {
			fmt.Println(err)
		}
		log.Println(string(requestDump))
		
		// pick params from header
		authorization := ctx.GetHeader("Authorization")
		device_id := ctx.GetHeader("device_id")
		signature_str := ctx.GetHeader("X-Birdol-Signature")
		timestamp := ctx.GetHeader("X-Birdol-TimeStamp")
		// Verify Authorization Header
		reg := regexp.MustCompile(`Bearer (.+)$`)
		if !reg.MatchString(authorization) {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrAuthorizationFail)
			ctx.Abort()
			return
		}
		access_token := reg.FindStringSubmatch(authorization)[1]

		// confirm accesstoken
		var recv_token model.AccessToken
		if err := database.Sqldb.Where("token = ?", access_token).First(&recv_token).Error; err != nil {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
			ctx.Abort()
			return
		}

		// confirm device id
		log.Println("recieved deviceid: ", device_id) // Added 3/21
		log.Println("token binded deviceid: ", recv_token.DeviceID) // Added 3/21
		
		if device_id != recv_token.DeviceID {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidDevice)
			ctx.Abort()
			return
		}

		// Verify signature
		encoded_xml_key := recv_token.PublicKey
		public_key, err := parseXML(encoded_xml_key)
		if err != nil {
			response.SetErrorResponse(ctx, http.StatusInternalServerError, response.ErrFailParseXML)
			ctx.Abort()
			return
		}
		body_byte, _ := io.ReadAll(ctx.Request.Body)
		ctx.Set("body_rawbyte", body_byte)
		request_body := string(body_byte)

		prefix := os.Getenv("API_VERSION")
		signature_base := prefix + ":" + timestamp + ":" + request_body
		hashed_base := sha256.Sum256([]byte(signature_base))
		signature, _ := hex.DecodeString(signature_str)

		verify_err := rsa.VerifyPKCS1v15(&public_key, crypto.SHA256, hashed_base[:], signature)
		if verify_err != nil {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidSignature)
			ctx.Abort()
			return
		}

		ctx.Set("access_token", recv_token)

		ctx.Next()
	}
}

/*
	Private functions
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
func parseXML(str string) (rsa.PublicKey, error) {
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
