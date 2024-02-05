package flygoose

import (
	"flygoose/pkg/beans"
	"flygoose/pkg/models"
	"flygoose/web/controllers/comm"
	"flygoose/web/services"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type SpecialController struct {
	comm.BaseComponent
}

func NewSpecialController() *SpecialController {
	return &SpecialController{}
}

func (c *SpecialController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/getSpecialList", "GetSpecialList")     //获取专栏列表
	b.Handle("POST", "/getSpecialDetail", "GetSpecialDetail") //获取专栏的详情
	b.Handle("POST", "/getSectionDetail", "GetSectionDetail") //获取小节的详情
}

func (c *SpecialController) GetSectionDetail() {
	var param beans.GetSectionParam
	if err := c.Ctx.ReadForm(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}
	srv := services.NewSpecialService()
	section, err := srv.GetSectionDetail(param.SectionId)
	if err != nil {
		c.RespFailedMessage("获取失败：" + err.Error())
		return
	}

	if section != nil && section.Status != models.SectionStatusPublished {
		c.RespFailedMessage("此章节未发布")
		return
	}

	c.RespSuccess(section, "获取小节成功")
}

func (c *SpecialController) GetSpecialDetail() {
	specialId, exist := c.PostFormInt64Default("specialId")
	if !exist || specialId == 0 {
		c.RespFailedMessage("参数错误")
		return
	}

	srv := services.NewSpecialService()
	special, list := srv.GetSpecialDetail(specialId)
	c.RespSuccess(iris.Map{
		"special": special,
		"list":    list,
	}, "获取成功")
}

func (c *SpecialController) GetSpecialList() {
	type Param struct {
		PageNum  int `json:"pageNum" validate:"required"`  //第几页，从1开始
		PageSize int `json:"pageSize" validate:"required"` //每页的数量
	}

	var param Param
	if err := c.Ctx.ReadForm(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	srv := services.NewSpecialService()
	list, hasMore := srv.GetSpecialList(param.PageNum, param.PageSize)
	c.RespSuccess(iris.Map{
		"list":    list,
		"hasMore": hasMore,
	}, "获取成功")
}
