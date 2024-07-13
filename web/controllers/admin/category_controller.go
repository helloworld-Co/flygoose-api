package admin

import (
	"flygoose/pkg/beans"
	"flygoose/pkg/models"
	"flygoose/web/controllers/comm"
	"flygoose/web/middlers"
	"flygoose/web/services"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"time"
)

type CategoryController struct {
	comm.BaseComponent
	CategorySrv *services.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{CategorySrv: services.NewCategoryService()}
}

func (c *CategoryController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/create", "Create", middlers.CheckAdminToken)                   //创建分类
	b.Handle("POST", "/update", "Update", middlers.CheckAdminToken)                   //更新分类
	b.Handle("POST", "/getCategoryList", "GetCategoryList", middlers.CheckAdminToken) //获取分类列表
}

func (c *CategoryController) GetCategoryList() {
	var bean beans.CatePageStatusBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	list := c.CategorySrv.GetCategoryList(bean.Status)
	c.RespSuccess(iris.Map{
		"list": list,
	}, "获取成功")
}

func (c *CategoryController) Create() {
	cate := models.Category{}
	err := c.Ctx.ReadForm(&cate)
	if err != nil {
		c.RespFailedMessage("参数错误")
		return
	}

	if cate.Name == "" {
		c.RespFailedMessage("分类名称不能为空")
		return
	}

	if cate.Status != models.CateStatusNormal && cate.Status != models.CateStatusOffline {
		c.RespFailedMessage("分类状态参数出错")
		return
	}

	cate.CreateTime = time.Now()
	cate.UpdateTime = time.Now()

	err = c.CategorySrv.Create(&cate)
	if err != nil {
		c.RespFailedMessage("创建分类失败")
	} else {
		c.RespSuccess(cate, "创建分类成功")
	}
}

func (c *CategoryController) Update() {
	fields := make([]string, 0)
	cate := models.Category{}

	id, _ := c.PostFormInt64("id")
	if id <= 0 {
		c.RespFailedMessage("id参数错误")
		return
	}

	name, exist := c.PostFormStringDefault("name")
	if exist {
		if name == "" {
			c.RespFailedMessage("分类名称不能为空")
			return
		} else {
			fields = append(fields, "Name")
			cate.Name = name
		}
	}

	seq, exist := c.PostFormIntDefault("seq")
	if exist {
		fields = append(fields, "Seq")
		cate.Seq = seq
	}

	icon, exist := c.PostFormStringDefault("icon")
	if exist {
		fields = append(fields, "Icon")
		cate.Icon = icon
	}

	font, exist := c.PostFormStringDefault("font")
	if exist {
		fields = append(fields, "Font")
		cate.Font = font
	}

	color, exist := c.PostFormStringDefault("color")
	if exist {
		fields = append(fields, "Color")
		cate.Color = color
	}

	status, exist := c.PostFormIntDefault("status")
	if exist {
		if status != 0 && status != 1 {
			c.RespFailedMessage("status参数错误")
			return
		}

		fields = append(fields, "Status")
		cate.Status = models.CateStatus(status)
	}

	if len(fields) > 0 {
		fields = append(fields, "UpdateTime")
		cate.UpdateTime = time.Now()
	}

	err := c.CategorySrv.Update(id, fields, &cate)
	if err != nil {
		c.RespFailedMessage("更新失败")
	} else {
		c.RespSuccess(cate, "更新成功")
	}
}
