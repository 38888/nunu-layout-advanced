package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	v1 "mine-admin/api/v1"
	"mine-admin/pkg/casbin"
	"mine-admin/pkg/ctxdata"
	"mine-admin/pkg/log"
	"net/http"
)

func StrictCasbin(r casbin.RBAC, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		roleId := ctxdata.GetRoleIdFromCtx(ctx)

		check, err := r.Check(roleId, method, path)
		if err != nil {
			v1.HandleError(ctx, http.StatusBadRequest, v1.ErAuthERR, nil)
			logger.WithContext(ctx).Error("Casbin error", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
				"roleId": roleId,
			}), zap.Error(err))
			ctx.Abort()
			return
		}
		if !check {
			v1.HandleError(ctx, http.StatusBadRequest, v1.ErAuthERR, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
