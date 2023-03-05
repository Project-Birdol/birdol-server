/*
Security middlewares for gin
*/
package security

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"regexp"

	"github.com/Project-Birdol/birdol-server/controller/jsonmodel"
	"github.com/Project-Birdol/birdol-server/database"
	"github.com/Project-Birdol/birdol-server/database/model"
	"github.com/Project-Birdol/birdol-server/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/ecdsa"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/publickey"
	"github.com/starkbank/ecdsa-go/v2/ellipticcurve/signature"
)

// Keytype
const (
	Rsa1024 = "rsa-1024"
	Rsa2048 = "rsa-2048"
	Rsa4096 = "rsa-4096"
	Ecdsa   = "ecdsa"
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
			// For backward compatibility, treat empty keytype as rsa1024
			if signatureAlgo != Rsa1024 || recvToken.KeyType != "" {
				response.SetErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidSignature)
				ctx.Abort()
				return
			} else {
				signatureAlgo = Rsa1024
			}
		}

		// Create message for verification
		bodyByte, _ := io.ReadAll(ctx.Request.Body)
		ctx.Set("body_rawbyte", bodyByte)
		requestBody := string(bodyByte)
		prefix := os.Getenv("API_VERSION")
		msg := prefix + ":" + timestamp + ":" + requestBody

		// Verify
		var verifyErr error
		switch signatureAlgo {
		case Rsa1024:
		case Rsa2048:
		case Rsa4096:
			verifyErr = verifyRsaSignature(msg, signatureStr, recvToken.PublicKey) 
		case Ecdsa:
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

// Inspect Publickey before registration
func InspectPublicKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.SetPrefix("[InspectPublicKey]")
		// Content Type
		contentType := ctx.GetHeader("Content-Type")
		if contentType != gin.MIMEJSON {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidType)
			ctx.Abort()
			return
		}

		// Extract key data from request body
		var keyInfo jsonmodel.KeyInspectInfo
		if err := ctx.ShouldBindBodyWith(&keyInfo, binding.JSON); err != nil {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrFailParseJSON)
			ctx.Abort()
			return
		}

		// Only ECDSA publickey is accepted for new device
		if keyInfo.KeyType != Ecdsa {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidKeyType)
			ctx.Abort()
			return
		}

		// Import test
		pubKeyStr, err := base64Decode(keyInfo.PublicKey)
		if err != nil {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidKey)
			ctx.Abort()
			return
		}
		pubKeyBlob, err := hex.DecodeString(string(pubKeyStr))
		if err != nil {
			response.SetErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidKey)
			ctx.Abort()
			return
		}

		_ = publickey.FromDer(pubKeyBlob) // panics when invalid key is passed

		log.Println("inspection passed.")
		ctx.Next()
	}
}

/*
	Private functions
*/

// Verify message signed with ECDSA Privatekey
func verifyEcdsaSignature(msg string, sigStr string, pubKeyStr string) error {
	pubKeyHexStr, err := base64Decode(pubKeyStr)
	if err != nil {
		return err
	}

	pubKeyBlob, err := hex.DecodeString(string(pubKeyHexStr))
	if err != nil {
		return err
	}

	pubkey := publickey.FromDer(pubKeyBlob)

	signatureHexStr, err := base64Decode(sigStr) 
	if err != nil {
		return err
	}

	sigbyte, err := hex.DecodeString(string(signatureHexStr))
	if err != nil {
		return err
	}

	signature := signature.FromDer(sigbyte)

	verified := ecdsa.Verify(msg, signature, &pubkey)
	if !verified {
		return errors.New("invalid signature passed")
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
