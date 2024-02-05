package flygoose

import (
	"flygoose/pkg/models"
	"flygoose/web/controllers/comm"
	"flygoose/web/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type SiteController struct {
	comm.BaseComponent
}

func NewSiteController() *SiteController {
	return &SiteController{}
}

func (c *SiteController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/getSiteInfo", "GetSiteInfo")             //获取网站信息
	b.Handle("POST", "/getWebmasterInfo", "GetWebmasterInfo")   //获取站长信息
	b.Handle("POST", "/getFriendLinkList", "GetFriendLinkList") //获取友链列表
	b.Handle("POST", "/getNoticeList", "GetNoticeList")         //获取通知列表
	b.Handle("POST", "/getBannerList", "GetBannerList")         //获取轮播图列表
}

func (c *SiteController) GetBannerList() {
	srv := services.NewBannerService()
	list, count := srv.GetBannerList(int(models.LinkStatusNormal), 1, 99999999)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取成功")
}

func (c *SiteController) GetNoticeList() {
	srv := services.NewNoticeService()
	list, count := srv.GetNoticeList(int(models.LinkStatusNormal), 1, 99999999)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取成功")
}

func (c *SiteController) GetFriendLinkList() {
	srv := services.NewLinkService()
	list, count := srv.GetLinkList(int(models.LinkStatusNormal), 1, 99999999)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取成功")
}

func (c *SiteController) GetWebmasterInfo() {
	srv := services.NewSiteService()
	webmaster := srv.GetWebmasterInfo()
	c.RespSuccess(webmaster, "获取成功")
}

func (c *SiteController) GetSiteInfo() {
	srv := services.NewSiteService()
	siteInfo := srv.GetUsedSiteInfo()
	c.RespSuccess(siteInfo, "获取成功")
}
