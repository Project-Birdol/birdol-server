package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestValidation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		signature_key := "test key"
		timestamp := ctx.GetHeader("X-Birdol-Request-Timestamp")
		signature := ctx.GetHeader("X-Birdol-Signature")

		buf := make([]byte, 2048)
		n, err := ctx.Request.Body.Read(buf)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです。",
			})
			ctx.Abort()
		}
		request_body := string(buf[0:n])
		ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(buf))

		signature_base := "v2:" + timestamp + ":" + request_body

		generated_signature := generateHmac(signature_base, signature_key)

		if signature != generated_signature {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "failed",
				"error":  "不適切なリクエストです。",
			})
			ctx.Abort()
		}

		ctx.Next()
	}
}

func generateHmac(base, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(base))
	return hex.EncodeToString(h.Sum(nil))
}