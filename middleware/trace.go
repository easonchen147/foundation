package middleware

import (
	"github.com/easonchen147/foundation/constant"
	"github.com/easonchen147/foundation/util"
	"github.com/gin-gonic/gin"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constant.TraceIdKey, util.GetNanoId())
	}
}
