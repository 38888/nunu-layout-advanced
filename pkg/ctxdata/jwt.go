package ctxdata

import (
	"github.com/38888/nunu-layout-advanced/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetUserIdFromCtx(ctx *gin.Context) int64 {
	v, exists := ctx.Get("claims")
	if !exists {
		return 0
	}
	uid, err := strconv.ParseInt(v.(*jwt.MyCustomClaims).UserId, 10, 64)
	if err != nil {
		return 0
	}

	return uid
}

func GetRoleIdFromCtx(ctx *gin.Context) int64 {
	v, exists := ctx.Get("claims")
	if !exists {
		return 0
	}
	uid, err := strconv.ParseInt(v.(*jwt.MyCustomClaims).RoleId, 10, 64)
	if err != nil {
		return 0
	}

	return uid
}
