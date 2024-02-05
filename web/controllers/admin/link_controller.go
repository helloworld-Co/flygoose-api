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

type LinkController struct {
	comm.BaseComponent
	LinkSrv *services.LinkService
}

func NewLinkController() *LinkController {
	return &LinkController{LinkSrv: services.NewLinkService()}
}

func (c *LinkController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/create", "Create", middlers.CheckAdminToken)           //创建
	b.Handle("POST", "/update", "Update", middlers.CheckAdminToken)           //更新
	b.Handle("POST", "/getLinkInfo", "GetLinkInfo", middlers.CheckAdminToken) //获取友链信息
	b.Handle("POST", "/getLinkList", "GetLinkList", middlers.CheckAdminToken) //获取友链
}

func (c *LinkController) GetLinkList() {
	var bean beans.LinkPageStatusBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	list, count := c.LinkSrv.GetLinkList(bean.Status, bean.PageNum, bean.PageSize)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取成功")
}

func (c *LinkController) GetLinkInfo() {
	id, exist := c.PostFormInt64Default("id")
	if !exist || id <= 0 {
		c.RespFailedMessage("id参数不存在或者错误")
		return
	}

	link := c.LinkSrv.GetLinkInfo(id)
	if link == nil {
		c.RespFailedMessage("友链不存在")
	} else {
		c.RespSuccess(link, "获取成功")
	}
}

func (c *LinkController) Update() {
	fields := make([]string, 0)
	link := &models.Link{}

	id, exist := c.PostFormInt64Default("id")
	if !exist || id <= 0 {
		c.RespFailedMessage("id参数不存在或者错误")
		return
	}

	title := c.PostFormString("title")
	if title != "" {
		fields = append(fields, "Title")
		link.Title = title
	}

	url := c.PostFormString("url")
	if url != "" {
		fields = append(fields, "Url")
		link.Url = url
	}

	seq, exist := c.PostFormIntDefault("seq")
	if exist {
		fields = append(fields, "Seq")
		link.Seq = seq
	}

	remark, exist := c.PostFormStringDefault("remark")
	if exist {
		fields = append(fields, "Remark")
		link.Remark = remark
	}

	validTime, exist := c.PostFormInt64Default("validTime")
	if exist {

		if validTime == 0 {
			//这是前段没有选择日期时给的值 默认给到2099-12-31 23:59:59
			fields = append(fields, "ValidTime")
			link.ValidTime = time.Date(2099, 12, 31, 23, 59, 59, 0, time.Local)
		} else {
			fields = append(fields, "ValidTime")
			link.ValidTime = time.Unix(validTime, 0)
		}
	}

	status, exist := c.PostFormIntDefault("status")
	if exist && (status == 0 || status == 1) {
		fields = append(fields, "Status")
		link.Status = models.LinkStatus(status)
	}

	if len(fields) == 0 {
		c.RespFailedMessage("请输入待修改的字段")
		return
	}

	fields = append(fields, "UpdateTime")
	link.UpdateTime = time.Now()

	err := c.LinkSrv.Update(id, fields, link)
	if err != nil {
		c.RespFailedMessage("更新失败")
	} else {
		c.RespSuccessMessage("更新成功")
	}
}

func (c *LinkController) Create() {
	var bean beans.CreateLinkBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	var validTime = time.Date(2099, 12, 31, 23, 59, 59, 0, time.Local)
	if bean.ValidTime != 0 {
		validTime = time.Unix(bean.ValidTime, 0)
	}

	link := models.Link{
		Title:      bean.Title,
		Url:        bean.Url,
		Seq:        bean.Seq,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		ValidTime:  validTime,
		Remark:     bean.Remark,
	}
	err := c.LinkSrv.Create(&link)
	if err != nil {
		c.RespFailedMessage("添加失败")
	} else {
		c.RespSuccessMessage("添加成功")
	}
}
