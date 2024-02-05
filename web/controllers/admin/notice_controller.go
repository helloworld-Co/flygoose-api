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

type NoticeController struct {
	comm.BaseComponent
	NoticeSrv *services.NoticeService
}

func NewNoticeController() *NoticeController {
	return &NoticeController{NoticeSrv: services.NewNoticeService()}
}

func (c *NoticeController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/create", "Create", middlers.CheckAdminToken)               //创建
	b.Handle("POST", "/update", "Update", middlers.CheckAdminToken)               //更新
	b.Handle("POST", "/getNoticeList", "GetNoticeList", middlers.CheckAdminToken) //获取通知列表
	b.Handle("POST", "/getNoticeInfo", "GetNoticeInfo", middlers.CheckAdminToken) //获取通知详情
}

func (c *NoticeController) GetNoticeInfo() {
	id, exist := c.PostFormInt64Default("id")
	if !exist || id <= 0 {
		c.RespFailedMessage("id参数不存在或者错误")
		return
	}

	notice := c.NoticeSrv.GetNoticeInfo(id)
	if notice == nil {
		c.RespFailedMessage("通知不存在")
	} else {
		c.RespSuccess(notice, "获取成功")
	}
}

func (c *NoticeController) GetNoticeList() {
	var bean beans.PageStatusBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	list, count := c.NoticeSrv.GetNoticeList(bean.Status, bean.PageNum, bean.PageSize)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取成功")
}

func (c *NoticeController) Update() {
	id, _ := c.PostFormInt64("id")
	if id <= 0 {
		c.RespFailedMessage("id参数错误")
		return
	}

	fields := make([]string, 0)
	m := models.Notice{}

	title, exist := c.PostFormStringDefault("title")
	if exist {
		if title == "" {
			c.RespFailedMessage("公告标题不能为空")
			return
		}
		fields = append(fields, "Title")
		m.Title = title
	}

	content, exist := c.PostFormStringDefault("content")
	if exist {
		fields = append(fields, "Content")
		m.Content = content
	}

	validTime, exist := c.PostFormInt64Default("validTime")
	if exist {
		fields = append(fields, "ValidTime")
		m.ValidTime = time.Unix(validTime, 0)
	}

	status, exist := c.PostFormIntDefault("status")
	if exist {
		if status != 0 && status != 1 {
			c.RespFailedMessage("status参数错误")
			return
		}
		fields = append(fields, "Status")
		m.Status = status
	}

	if len(fields) == 0 {
		c.RespFailedMessage("请输入待更新的字段")
		return
	}

	err := c.NoticeSrv.Update(id, fields, &m)
	if err != nil {
		c.RespFailedMessage("更新失败")
	} else {
		c.RespSuccess(m, "更新成功")
	}
}

func (c *NoticeController) Create() {
	m := models.Notice{}

	title, _ := c.PostFormStringDefault("title")
	if title == "" {
		c.RespFailedMessage("公告标题不能为空")
		return
	}
	m.Title = title

	content, _ := c.PostFormStringDefault("content")
	m.Content = content

	validTime, exist := c.PostFormInt64Default("validTime")
	if exist {
		m.ValidTime = time.Unix(validTime, 0)
	} else {
		m.ValidTime = time.Date(2099, 12, 31, 23, 59, 59, 0, time.Local)
		validTime = m.ValidTime.Unix()
	}

	status, err := c.PostFormInt("status")
	if err != nil {
		c.RespFailedMessage("请检查状态")
		return
	}

	if status != 1 && status != 0 {
		c.RespFailedMessage("请检查状态参数")
		return
	}

	m.CreateTime = time.Now()
	m.UpdateTime = time.Now()
	m.Status = status

	err = c.NoticeSrv.Create(&m)
	if err != nil {
		c.RespFailedMessage("创建失败")
	} else {
		c.RespSuccess(m, "创建成功")
	}
}
