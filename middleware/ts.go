package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/easonchen147/foundation/cfg"

	"github.com/gin-gonic/gin"
)

const (
	HeaderTsKey = "Ts"
)

type TsVerify struct {
	Expire time.Duration
}

func (b *TsVerify) verify(c *gin.Context) (bool, error) {
	origin := c.Request.Header.Get(HeaderTsKey)
	second, err := strconv.ParseInt(origin, 10, 64)
	if err != nil {
		return false, err
	}
	originTime := time.Unix(second, 0)
	return time.Since(originTime) > b.Expire, nil
}

var tsVerifier *TsVerify

func init() {
	if cfg.AppConf.TsConfig == nil || cfg.AppConf.TsConfig.Expire == "" {
		return
	}
	expire, err := time.ParseDuration(cfg.AppConf.TsConfig.Expire)
	if err != nil {
		return
	}
	tsVerifier = &TsVerify{
		Expire: expire,
	}
}

func VerifyTs() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tsVerifier == nil { //未配置接口请求时间戳校验的直接跳过
			c.Next()
			return
		}

		ok, err := tsVerifier.verify(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
	}
}
