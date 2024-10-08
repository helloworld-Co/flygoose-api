package admin

import (
	"flygoose/web/controllers/comm"
	"flygoose/web/middlers"
	"flygoose/web/services"
	"github.com/kataras/iris/v12/mvc"
)

type WorkStationController struct {
	comm.BaseComponent
}

func NewWorkStationController() *WorkStationController {
	return &WorkStationController{}
}

func (c *WorkStationController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/getStatistics", "GetStatistics", middlers.CheckAdminToken) //获取统计数据
}

func (c *WorkStationController) GetStatistics() {

	c.RespSuccess(services.NewWorkStationService().GetStatistics(), "success")
}
