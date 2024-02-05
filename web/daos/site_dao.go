package daos

import (
	"flygoose/datasource"
	"flygoose/pkg/models"
	"flygoose/pkg/tlog"
	"gorm.io/gorm"
)

type SiteDao struct {
	db *gorm.DB
}

func NewSiteDao() *SiteDao {
	return &SiteDao{db: datasource.GetMasterDB()}
}

func (dao *SiteDao) CreateSite(site *models.Site) error {
	result := dao.db.Create(site)
	if result.Error != nil {
		tlog.Error2("SiteDao:CreateSite 出错", result.Error)
	}
	return result.Error
}

func (dao *SiteDao) UpdateSite(id int64, fields []string, site *models.Site) error {
	result := dao.db.Model(&models.Site{}).Where("id=?", id).Select(fields).Updates(site)
	if result.Error != nil {
		tlog.Error2("SiteDao:UpdateSite 出错", result.Error)
	}
	return result.Error
}

func (dao *SiteDao) GetSiteInfoList() []models.Site {
	var list []models.Site
	result := dao.db.Order("status desc , update_time desc").Find(&list)
	if result.Error != nil {
		tlog.Error2("SiteDao:GetSiteInfoList 出错", result.Error)
		return nil
	}

	return list
}

func (dao *SiteDao) GetUsedSiteInfo() *models.Site {
	var list []models.Site
	result := dao.db.Model(&models.Site{}).Where("status=?", models.SiteStatusOnline).Order("status desc , id desc").Limit(1).Offset(0).Find(&list)
	if result.Error != nil {
		tlog.Error2("SiteDao:GetUsedSiteInfo 出错", result.Error)
		return nil
	}
	if len(list) == 0 {
		return nil
	}

	return &list[0]
}

func (dao *SiteDao) GetWebmasterInfo() *models.Webmaster {
	var webmaster models.Webmaster
	result := dao.db.Model(&models.Webmaster{}).Order("id").Limit(1).First(&webmaster)
	if result.Error != nil {
		tlog.Error2("SiteDao:GetWebmasterInfo 出错", result.Error)
		return nil
	}
	return &webmaster
}

func (dao *SiteDao) CreateWebmaster(webmaster *models.Webmaster) error {
	result := dao.db.Create(webmaster)
	if result.Error != nil {
		tlog.Error2("SiteDao:CreateWebmaster 出错", result.Error)
	}
	return result.Error
}

func (dao *SiteDao) UpdateWebmaster(id int64, fields []string, webmaster *models.Webmaster) error {
	result := dao.db.Model(&models.Webmaster{}).Where("id=?", id).Select(fields).Updates(webmaster)
	if result.Error != nil {
		tlog.Error2("SiteDao:UpdateWebmaster 出错", result.Error)
	}
	return result.Error
}

func (dao *SiteDao) GetWebmasterInfoList() []models.Webmaster {
	var list []models.Webmaster
	result := dao.db.Model(&models.Webmaster{}).Order("status desc , update_time desc").Find(&list)
	if result.Error != nil {
		tlog.Error2("SiteDao:GetWebmasterInfoList 出错", result.Error)
		return nil
	}
	return list
}
