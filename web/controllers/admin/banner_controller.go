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

type BannerController struct {
	comm.BaseComponent
}

func NewBannerController() *BannerController {
	return &BannerController{}
}

func (c *BannerController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/create", "Create", middlers.CheckAdminToken)               //创建
	b.Handle("POST", "/update", "Update", middlers.CheckAdminToken)               //更新
	b.Handle("POST", "/getBannerInfo", "GetBannerInfo", middlers.CheckAdminToken) //获取banner信息
	b.Handle("POST", "/getBannerList", "GetBannerList", middlers.CheckAdminToken) //获取banner列表
}

func (c *BannerController) GetBannerList() {
	var bean beans.BannerPageStatusBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	list, count := services.NewBannerService().GetBannerList(bean.Status, bean.PageNum, bean.PageSize)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取成功")
}

func (c *BannerController) GetBannerInfo() {
	id, exist := c.PostFormInt64Default("id")
	if !exist || id <= 0 {
		c.RespFailedMessage("id参数不存在或者错误")
		return
	}

	banner := services.NewBannerService().GetBannerInfo(id)
	if banner == nil {
		c.RespFailedMessage("banner不存在")
	} else {
		c.RespSuccess(banner, "获取成功")
	}
}

func (c *BannerController) Update() {
	fields := make([]string, 0)
	banner := &models.Banner{}

	id, exist := c.PostFormInt64Default("id")
	if !exist || id <= 0 {
		c.RespFailedMessage("id参数不存在或者错误")
		return
	}

	title := c.PostFormString("title")
	if title != "" {
		fields = append(fields, "Title")
		banner.Title = title
	}

	url := c.PostFormString("url")
	if url != "" {
		fields = append(fields, "Url")
		banner.Url = url
	}

	targetUrl := c.PostFormString("targetUrl")
	if targetUrl != "" {
		fields = append(fields, "TargetUrl")
		banner.TargetUrl = targetUrl
	}

	seq, exist := c.PostFormIntDefault("seq")
	if exist {
		fields = append(fields, "Seq")
		banner.Seq = seq
	}

	status, exist := c.PostFormIntDefault("status")
	if exist && (status == 0 || status == 1) {
		fields = append(fields, "Status")
		banner.Status = models.BannerStatus(status)
	}

	if len(fields) == 0 {
		c.RespFailedMessage("请输入待修改的字段")
		return
	}

	fields = append(fields, "CreateTime")
	banner.CreateTime = time.Now()

	err := services.NewBannerService().Update(id, fields, banner)
	if err != nil {
		c.RespFailedMessage("更新失败:" + err.Error())
	} else {
		c.RespSuccessMessage("更新成功")
	}
}

func (c *BannerController) Create() {
	var bean beans.CreateBannerBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	banner := models.Banner{
		Title:      bean.Title,
		Url:        bean.Url,
		TargetUrl:  bean.TargetUrl,
		Seq:        bean.Seq,
		CreateTime: time.Now(),
	}

	err := services.NewBannerService().Create(&banner)
	if err != nil {
		c.RespFailedMessage("添加失败")
	} else {
		c.RespSuccessMessage("添加成功")
	}
}
