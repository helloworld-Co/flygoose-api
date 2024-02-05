package middlers

import (
	"flygoose/pkg/models"
	"flygoose/web/caching"
	"flygoose/web/controllers/comm"
	"flygoose/web/daos"
	"github.com/kataras/iris/v12"
	"time"
)

// CheckAdminToken 判断 huskar app token是否过期
func CheckAdminToken(ctx iris.Context) {
	token := ctx.Request().Header.Get("token")
	if token == "" {
		ctx.JSON(comm.GooseEgg{Code: comm.CodeTokenExpired, Message: "无效的凭证"})
		return
	}

	var admin *models.Admin
	accessDao := daos.NewAccessDao()
	uid := caching.GetAdminCache().GetUid(token)
	if uid <= 0 {
		admin = accessDao.FindAminWithToken(token)
		if admin != nil && admin.ValidTime.After(time.Now()) {
			caching.GetAdminCache().SetUid(token, admin.ID)
			ctx.Values().Set("user_id", admin.ID)
		} else {
			ctx.Values().Set("user_id", 0)
			ctx.JSON(comm.GooseEgg{Code: comm.CodeTokenExpired, Message: "无效凭证"})
			return
		}
	}

	if admin == nil {
		admin = accessDao.FindAdminByUid(uid)
		if admin == nil {
			ctx.JSON(comm.GooseEgg{Code: comm.CodeTokenExpired, Message: "无效凭证"})
			return
		}
	}

	if admin == nil {
		ctx.JSON(comm.GooseEgg{Code: comm.CodeTokenExpired, Message: "无效凭证"})
		return
	}

	if admin.Status == models.AdminStatusBlocked {
		caching.GetAdminCache().DeleteUid(token)
		accessDao.CleanUserLoginInfo(uid)
		ctx.JSON(comm.GooseEgg{Code: comm.CodeFailed, Message: "账号被封，请联系管理员"})
		return
	}

	if time.Now().After(admin.ValidTime) {
		caching.GetAdminCache().DeleteUid(token)
		ctx.JSON(comm.GooseEgg{Code: comm.CodeTokenExpired, Message: "会话已过期"})
		return
	}

	ctx.Values().Set("user_id", admin.ID)
	ctx.Next()
}
