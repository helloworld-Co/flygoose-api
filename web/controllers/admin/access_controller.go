package admin

import (
	"flygoose/web/controllers/comm"
	"flygoose/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AccessController struct {
	comm.BaseComponent
	AccessSrv *services.AccessService
}

func NewAccessController() *AccessController {
	return &AccessController{AccessSrv: services.NewAccessService()}
}

func (c *AccessController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/login", "Login")   //登录
	b.Handle("POST", "/logout", "Logout") //退出登录
}

func (c *AccessController) Login() {
	username, exist := c.PostFormStringDefault("username")
	if !exist || username == "" {
		c.RespFailedMessage("参数错误")
		return
	}

	password, exist := c.PostFormStringDefault("password")
	if !exist || password == "" {
		c.RespFailedMessage("参数错误")
		return
	}

	ok, str := c.AccessSrv.LoginIn(username, password)
	if !ok {
		c.RespFailedMessage("登录失败：" + str)
		return
	}

	c.RespSuccess(iris.Map{
		"token": str,
	}, "登录成功")
}

func (c *AccessController) Logout() {
	token := c.Ctx.Request().Header.Get("token")
	if token == "" {
		c.RespFailedMessage("token参数错误")
		return
	}

	admin, err := c.AccessSrv.FirstAdminByToken(token)
	if admin == nil || err != nil {
		c.RespFailedMessage("退出失败，无此管理账户")
		return
	}

	err = c.AccessSrv.Logout(admin.ID)
	if err != nil {
		c.RespFailedMessage(err.Error())
	} else {
		c.RespSuccessMessage("退出成功")
	}
}
