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

type SpecialController struct {
	comm.BaseComponent
	Srv *services.SpecialService
}

func NewSpecialController() *SpecialController {
	return &SpecialController{Srv: services.NewSpecialService()}
}

func (c *SpecialController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/create", "Create", middlers.CheckAdminToken)               //创建
	b.Handle("POST", "/update", "Update", middlers.CheckAdminToken)               //更新
	b.Handle("POST", "/searchSpecial", "SearchSpecial", middlers.CheckAdminToken) //搜索专栏

	b.Handle("POST", "/addSection", "AddSection", middlers.CheckAdminToken)             //添加小节
	b.Handle("POST", "/updateSection", "UpdateSection", middlers.CheckAdminToken)       //更新小节
	b.Handle("POST", "/getSectionList", "GetSectionList", middlers.CheckAdminToken)     //获取小节列表
	b.Handle("POST", "/getSectionDetail", "GetSectionDetail", middlers.CheckAdminToken) //获取小节详情
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

	section, err := c.Srv.GetSectionDetail(param.SectionId)
	if err != nil {
		c.RespFailedMessage("获取失败：" + err.Error())
		return
	}

	if section != nil && section.Status == models.SectionStatusDeleted {
		c.RespFailedMessage("小节不存在")
		return
	}

	c.RespSuccess(section, "获取小节详情成功")
}

func (c *SpecialController) GetSectionList() {
	specialId, exist := c.PostFormInt64Default("specialId")
	if !exist || specialId <= 0 {
		c.RespFailedMessage("specialId参数错误")
		return
	}

	//参数 models.SectionStatus, 其中status=0,是获取全部状态的小节
	status, err := c.PostFormInt("status")
	if err != nil {
		c.RespFailedMessage("status参数错误")
		return
	}

	list, count := c.Srv.GetSectionList(specialId, status)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取数据成功")
}

func (c *SpecialController) UpdateSection() {
	fields := make([]string, 0)
	section := &models.Section{}

	sectionId, exist := c.PostFormInt64Default("sectionId") //sectionId
	specialId, exist := c.PostFormInt64Default("specialId") //specialId

	title := c.PostFormString("title")
	if title != "" {
		fields = append(fields, "Title")
		section.Title = title
	}

	content, exist := c.PostFormStringDefault("content")
	if exist {
		fields = append(fields, "Content")
		section.Content = content
	}

	html, exist := c.PostFormStringDefault("html")
	if exist {
		fields = append(fields, "Html")
		section.Html = html
	}

	tags, exist := c.PostFormStringDefault("tags")
	if exist {
		fields = append(fields, "Tags")
		section.Tags = tags
	}

	isHtml, exist := c.PostFormIntDefault("isHtml")
	if exist && (isHtml == 0 || isHtml == 1) {
		fields = append(fields, "IsHtml")
		section.IsHtml = isHtml
	}

	seq, exist := c.PostFormIntDefault("seq")
	if exist {
		fields = append(fields, "Seq")
		section.Seq = seq
	}

	status, exist := c.PostFormIntDefault("status")
	if exist && (models.SectionStatus(status) == models.SectionStatusCreated ||
		models.SectionStatus(status) == models.SectionStatusDeleted ||
		models.SectionStatus(status) == models.SectionStatusPublished) {
		fields = append(fields, "Status")
		section.Status = models.SectionStatus(status)
	}

	if len(fields) == 0 {
		c.RespFailedMessage("请输入要修改的参数")
		return
	}

	section.UpdateTime = time.Now()
	fields = append(fields, "UpdateTime")

	if section.Status == models.SectionStatusPublished {
		section.PublishTime = time.Now()
		fields = append(fields, "PublishTime")
	}

	err := c.Srv.UpdateSection(specialId, sectionId, fields, section)
	if err != nil {
		c.RespFailedMessage("更新失败:" + err.Error())
	} else {
		c.RespSuccessMessage("更新成功")
	}
}

func (c *SpecialController) AddSection() {
	var param beans.AddSectionParam
	if err := c.Ctx.ReadForm(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	section, err := c.Srv.AddSection(&param)
	if err != nil {
		c.RespFailedMessage("添加失败:" + err.Error())
	} else {
		c.RespSuccess(section, "添加成功")
	}
}

func (c *SpecialController) SearchSpecial() {
	var param beans.SearchSpecialParam
	if err := c.Ctx.ReadForm(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	list, count := c.Srv.SearchSpecial(&param)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取数据成功")
}

func (c *SpecialController) Update() {
	fields := make([]string, 0)
	special := &models.Special{}

	id, exist := c.PostFormInt64Default("id")
	if !exist || id <= 0 {
		c.RespFailedMessage("id参数不存在或者错误")
		return
	}

	title := c.PostFormString("title")
	if title != "" {
		fields = append(fields, "Title")
		special.Title = title
	}

	intro, exist := c.PostFormStringDefault("intro")
	if exist {
		fields = append(fields, "Intro")
		special.Intro = intro
	}

	cover, exist := c.PostFormStringDefault("cover")
	if exist {
		fields = append(fields, "Cover")
		special.Cover = cover
	}

	status, exist := c.PostFormIntDefault("status")
	if exist {
		var ss = models.SpecialStatus(status)
		if ss != models.SpecialStatusPublished && ss != models.SpecialStatusDeleted && ss != models.SpecialStatusCreated {
			c.RespFailedMessage("状态参数错误")
			return
		}

		fields = append(fields, "Status")
		special.Status = ss
	}

	if len(fields) == 0 {
		c.RespFailedMessage("请输入要修改的参数")
		return
	}

	fields = append(fields, "UpdateTime")
	special.UpdateTime = time.Now()

	err := c.Srv.Update(id, fields, special)
	if err != nil {
		c.RespFailedMessage("更新失败")
	} else {
		c.RespSuccessMessage("更新成功")
	}
}

func (c *SpecialController) Create() {
	var param beans.CreateSpecialParam
	if err := c.Ctx.ReadForm(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	special, err := c.Srv.Create(&param)
	if err != nil {
		c.RespFailedMessage("创建失败")
	} else {
		c.RespSuccess(special, "创建成功")
	}
}
