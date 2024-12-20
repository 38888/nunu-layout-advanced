package handler

import (
	"github.com/38888/nunu-layout-advanced/pkg/jwt"
	"github.com/38888/nunu-layout-advanced/pkg/log"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(
	logger *log.Logger,
) *Handler {
	return &Handler{
		logger: logger,
	}
}
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
