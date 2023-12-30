package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/easonchen147/foundation/cfg"

	"github.com/gin-gonic/gin"
)

const (
	HeaderSignKey = "Sk"
)

type BodySigner struct {
	secret []byte
	salt   []byte
}

func (b *BodySigner) sign(c *gin.Context) (string, error) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", err
	}
	defer func() {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	}()

	m := hmac.New(sha256.New, b.secret)
	m.Write(append(data, b.salt...))
	return hex.EncodeToString(m.Sum(nil)), nil
}

func (b *BodySigner) verify(c *gin.Context) (bool, error) {
	origin := c.Request.Header.Get(HeaderSignKey)
	current, err := b.sign(c)
	if err != nil {
		return false, err
	}
	return origin == current, nil
}

var bodySigner *BodySigner

func init() {
	if cfg.AppConf.SignConfig == nil || cfg.AppConf.SignConfig.Secret == "" {
		return
	}
	bodySigner = &BodySigner{
		secret: []byte(cfg.AppConf.SignConfig.Secret),
		salt:   []byte(cfg.AppConf.SignConfig.Salt),
	}
}

func VerifySk() gin.HandlerFunc {
	return func(c *gin.Context) {
		if bodySigner == nil { //未配置签名校验秘钥的直接跳过
			c.Next()
			return
		}

		match, err := bodySigner.verify(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		if !match {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
	}
}
