/*
Middlewares for gin
*/
package middlewares

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io"
	"math/big"
	"net/http"
	"os"
	"regexp"

	"github.com/MISW/birdol-server/database"
	"github.com/MISW/birdol-server/database/model"
	"github.com/MISW/birdol-server/utils/response"
	"github.com/gin-gonic/gin"
)

// RSA PublicKey
type rsaPublicKey struct {
	Modulus	string
	Exponent string
}

/*
	Main function
*/

// Verify request using signature
func RequestValidation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// pick params from header
		authorization := ctx.GetHeader("Authorization")
		deviceId := ctx.GetHeader("DeviceID")
		signatureStr := ctx.GetHeader("X-Birdol-Signature")
		timestamp := ctx.GetHeader("X-Birdol-TimeStamp")
		signatureAlgo := ctx.GetHeader("X-Birdol-Signature-Algo");

		// Verify Authorization Header
		reg := regexp.MustCompile(`Bearer (.+)$`)
		if !reg.MatchString(authorization) {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrAuthorizationFail)
			ctx.Abort()
			return
		}
		accessToken := reg.FindStringSubmatch(authorization)[1]

		// confirm accesstoken
		var recvToken model.AccessToken
		if err := database.Sqldb.Where("token = ?", accessToken).First(&recvToken).Error; err != nil {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
			ctx.Abort()
			return
		}

		// Confirm device uuid
		if deviceId != recvToken.DeviceID {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidDevice)
			ctx.Abort()
			return
		}

		// Check signature algorithm
		if signatureAlgo != recvToken.KeyType {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidSignature)
			ctx.Abort()
			return
		}

		// Create message for verification
		bodyByte, _ := io.ReadAll(ctx.Request.Body)
		ctx.Set("bodyByte", bodyByte)
		requestBody := string(bodyByte)
		prefix := os.Getenv("API_VERSION")
		msg := prefix + ":" + timestamp + ":" + requestBody

		// Verify
		var verifyErr error
		switch signatureAlgo {
		case "rsa-1024":
		case "rsa-2048":
		case "rsa-4096":
			verifyErr = verifyRsaSignature(msg, signatureStr, recvToken.PublicKey) 
		case "ecdsa":
			verifyErr = verifyEcdsaSignature(msg, signatureStr, recvToken.PublicKey)
		default:
			verifyErr = errors.New("invalid signature algorithm")
		}

		if verifyErr != nil {
			response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidSignature)
			ctx.Abort()
			return
		}

		// Set AccessToken struct to ctx
		ctx.Set("access_token", recvToken)

		ctx.Next()
	}
}

/*
	Private functions
*/

// Verify message signed with ECDSA Privatekey
func verifyEcdsaSignature(msg string, sigStr string, pubKeyStr string) error {
	pubKeyBlob, err := base64Decode(pubKeyStr)
	if err != nil {
		return err
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBlob)
	if err != nil {
		return err
	}

	signature, err := base64Decode(sigStr) 
	if err != nil {
		return err
	}

	hashedMsg := sha512.Sum512([]byte(msg))

	if !ecdsa.VerifyASN1(pubKey.(*ecdsa.PublicKey), hashedMsg[:], signature) {
		return errors.New("invalid signature found")
	}

	return nil
}

// Verify message signed with RSA Privatekey
func verifyRsaSignature(msg string, sigStr string, pubKeyStr string) error {
	rsaPubKey, err := parseXML(pubKeyStr)
	if err != nil {
		return err
	}

	signature, err := hex.DecodeString(sigStr)
	if err != nil {
		return err
	}

	hashedMsg := sha512.Sum512([]byte(msg))

	return rsa.VerifyPKCS1v15(&rsaPubKey, crypto.SHA512, hashedMsg[:], signature)
}

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
