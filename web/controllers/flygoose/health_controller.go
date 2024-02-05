package flygoose

import (
	"flygoose/web/controllers/comm"
	"github.com/kataras/iris/v12/mvc"
)

type HealthController struct {
	comm.BaseComponent
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "Health") //获取博客分类列表
}

func (c *HealthController) Health() {

	c.RespSuccessMessage("success")
}
