package admin

import (
	"flygoose/pkg/models"
	"flygoose/web/controllers/comm"
	"flygoose/web/middlers"
	"flygoose/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"time"
)

type SiteController struct {
	comm.BaseComponent
	SiteSrv *services.SiteService
}

func NewSiteController() *SiteController {
	return &SiteController{SiteSrv: services.NewSiteService()}
}

func (c *SiteController) BeforeActivation(b mvc.BeforeActivation) {
	//网站信息
	b.Handle("POST", "/createSite", "CreateSite", middlers.CheckAdminToken)           //创建
	b.Handle("POST", "/updateSite", "UpdateSite", middlers.CheckAdminToken)           //更新
	b.Handle("POST", "/getSiteInfoList", "GetSiteInfoList", middlers.CheckAdminToken) //获取网站信息
	b.Handle("POST", "/getUsedSiteInfo", "GetUsedSiteInfo", middlers.CheckAdminToken) //获取在网站上展示的那一条

	//站长信息
	b.Handle("POST", "/createWebmasterInfo", "CreateWebmasterInfo", middlers.CheckAdminToken)   //创建站长信息
	b.Handle("POST", "/updateWebmasterInfo", "UpdateWebmasterInfo", middlers.CheckAdminToken)   //更新站长信息
	b.Handle("POST", "/getWebmasterInfo", "GetWebmasterInfo", middlers.CheckAdminToken)         //获取站长信息
	b.Handle("POST", "/getWebmasterInfoList", "GetWebmasterInfoList", middlers.CheckAdminToken) //获取站长信息
}

func (c *SiteController) CreateWebmasterInfo() {
	nicker := c.PostFormString("nicker")
	intro := c.PostFormString("intro")
	slogan := c.PostFormString("slogan")
	avatar := c.PostFormString("avatar")
	job := c.PostFormString("job")
	email := c.PostFormString("email")
	qq := c.PostFormString("qq")
	wechat := c.PostFormString("wechat")
	rewardCode := c.PostFormString("rewardCode")

	webmaster := models.Webmaster{
		Nicker:     nicker,
		Slogan:     slogan,
		Intro:      intro,
		Avatar:     avatar,
		Job:        job,
		Email:      email,
		QQ:         qq,
		Wechat:     wechat,
		RewardCode: rewardCode,
		UpdateTime: time.Now(),
	}

	err := c.SiteSrv.CreateWebmaster(&webmaster)
	if err != nil {
		c.RespFailedMessage("创建失败")
	} else {
		c.RespSuccess(webmaster, "创建成功")
	}
}

func (c *SiteController) UpdateWebmasterInfo() {
	fields := make([]string, 0)
	webmaster := &models.Webmaster{}

	id, exist := c.PostFormInt64Default("id")
	if !exist || id <= 0 {
		c.RespFailedMessage("id参数不存在或者错误")
		return
	}

	intro, exist := c.PostFormStringDefault("intro")
	if exist {
		fields = append(fields, "Intro")
		webmaster.Intro = intro
	}

	slogan, exist := c.PostFormStringDefault("slogan")
	if exist {
		fields = append(fields, "Slogan")
		webmaster.Slogan = slogan
	}

	nicker, exist := c.PostFormStringDefault("nicker")
	if exist {
		fields = append(fields, "Nicker")
		webmaster.Nicker = nicker
	}

	avatar, exist := c.PostFormStringDefault("avatar")
	if exist {
		fields = append(fields, "Avatar")
		webmaster.Avatar = avatar
	}

	job, exist := c.PostFormStringDefault("job")
	if exist {
		fields = append(fields, "Job")
		webmaster.Job = job
	}

	email, exist := c.PostFormStringDefault("email")
	if exist {
		fields = append(fields, "Email")
		webmaster.Email = email
	}

	qq, exist := c.PostFormStringDefault("qq")
	if exist {
		fields = append(fields, "QQ")
		webmaster.QQ = qq
	}

	wechat, exist := c.PostFormStringDefault("wechat")
	if exist {
		fields = append(fields, "Wechat")
		webmaster.Wechat = wechat
	}

	rewardCode, exist := c.PostFormStringDefault("rewardCode")
	if exist {
		fields = append(fields, "RewardCode")
		webmaster.RewardCode = rewardCode
	}

	status, exist := c.PostFormIntDefault("status")
	if exist && (models.WebmasterStatus(status) == models.SiteStatusOnline || models.WebmasterStatus(status) == models.SiteStatusOffline) {
		fields = append(fields, "Status")
		webmaster.Status = models.WebmasterStatus(status)
	}

	if len(fields) == 0 {
		c.RespFailedMessage("请输入待更新的字段")
		return
	}

	fields = append(fields, "UpdateTime")
	webmaster.UpdateTime = time.Now()

	err := c.SiteSrv.UpdateWebmaster(id, fields, webmaster)
	if err != nil {
		c.RespFailedMessage("更新失败")
	} else {
		c.RespSuccess(webmaster, "更新成功")
	}
}

func (c *SiteController) GetWebmasterInfoList() {
	webmasters := c.SiteSrv.GetWebmasterInfoList()
	if webmasters == nil {
		c.RespFailedMessage("还未创建站长信息")
	} else {
		c.RespSuccess(iris.Map{
			"list": webmasters,
		}, "获取数据成功")
	}
}

func (c *SiteController) GetWebmasterInfo() {
	webmaster := c.SiteSrv.GetWebmasterInfo()
	if webmaster == nil {
		c.RespFailedMessage("还未创建站长信息")
	} else {
		c.RespSuccess(webmaster, "获取数据成功")
	}
}

func (c *SiteController) GetUsedSiteInfo() {
	site := c.SiteSrv.GetUsedSiteInfo()
	c.RespSuccess(site, "获取数据成功")
}

func (c *SiteController) GetSiteInfoList() {
	//mengxianhou@20240122 这里只返回1个
	firstSite := c.SiteSrv.GetUsedSiteInfo()

	result := make([]models.Site, 1)
	result[0] = *firstSite
	c.RespSuccess(iris.Map{
		"list": result,
	}, "获取数据成功")
}

func (c *SiteController) UpdateSite() {
	fields := make([]string, 0)
	site := &models.Site{}

	id, exist := c.PostFormInt64Default("id")
	if !exist || id <= 0 {
		c.RespFailedMessage("id参数不存在或者错误")
		return
	}

	title, exist := c.PostFormStringDefault("title")
	if exist {
		fields = append(fields, "Title")
		site.Title = title
	}

	intro, exist := c.PostFormStringDefault("intro")
	if exist {
		fields = append(fields, "Intro")
		site.Intro = intro
	}

	slogan, exist := c.PostFormStringDefault("slogan")
	if exist {
		fields = append(fields, "Slogan")
		site.Slogan = slogan
	}

	cover, exist := c.PostFormStringDefault("cover")
	if exist {
		fields = append(fields, "Cover")
		site.Cover = cover
	}

	copyright, exist := c.PostFormStringDefault("copyright")
	if exist {
		fields = append(fields, "Copyright")
		site.Copyright = copyright
	}

	icp, exist := c.PostFormStringDefault("icp")
	if exist {
		fields = append(fields, "Icp")
		site.Icp = icp
	}

	extra, exist := c.PostFormStringDefault("extra")
	if exist {
		fields = append(fields, "Extra")
		site.Extra = extra
	}

	status, exist := c.PostFormIntDefault("status")
	if exist && (models.SiteStatus(status) == models.SiteStatusOnline || models.SiteStatus(status) == models.SiteStatusOffline) {
		fields = append(fields, "Status")
		site.Status = models.SiteStatus(status)
	}

	if len(fields) == 0 {
		c.RespFailedMessage("请输入待更新的字段")
		return
	}

	fields = append(fields, "UpdateTime")
	site.UpdateTime = time.Now()

	err := c.SiteSrv.UpdateSite(id, fields, site)
	if err != nil {
		c.RespFailedMessage("更新失败")
	} else {
		c.RespSuccess(site, "更新成功")
	}
}

func (c *SiteController) CreateSite() {
	title := c.PostFormString("title")
	if title == "" {
		title = "我的站点"
	}

	intro := c.PostFormString("intro")
	if intro == "" {
		intro = "helloworld开发者社区旗下开源博客系统"
	}

	slogan := c.PostFormString("slogan")
	if slogan == "" {
		slogan = "同一个世界，同一行代码"
	}

	cover := c.PostFormString("cover")

	copyright := c.PostFormString("copyright")
	if copyright == "" {
		copyright = "本站由helloworld.net开发者社区提供技术支持"
	}

	icp := c.PostFormString("icp")
	extra := c.PostFormString("extra")

	site := models.Site{
		Title:      title,
		Intro:      intro,
		Slogan:     slogan,
		Cover:      cover,
		Copyright:  copyright,
		Icp:        icp,
		Extra:      extra,
		Status:     models.CateStatusOffline,
		UpdateTime: time.Now(),
	}

	err := c.SiteSrv.CreateSite(&site)
	if err != nil {
		c.RespFailedMessage("创建失败")
	} else {
		c.RespSuccess(site, "创建成功")
	}
}
