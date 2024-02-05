package services

import (
	"flygoose/pkg/models"
	"flygoose/web/daos"
)

type SiteService struct {
	siteDao *daos.SiteDao
}

func NewSiteService() *SiteService {
	return &SiteService{siteDao: daos.NewSiteDao()}
}

func (s *SiteService) CreateSite(site *models.Site) error {
	return s.siteDao.CreateSite(site)
}

func (s *SiteService) UpdateSite(id int64, fields []string, site *models.Site) error {
	return s.siteDao.UpdateSite(id, fields, site)
}

func (s *SiteService) GetSiteInfoList() []models.Site {
	return s.siteDao.GetSiteInfoList()
}

func (s *SiteService) GetUsedSiteInfo() *models.Site {
	return s.siteDao.GetUsedSiteInfo()
}

func (s *SiteService) GetWebmasterInfo() *models.Webmaster {
	return s.siteDao.GetWebmasterInfo()
}

func (s *SiteService) GetWebmasterInfoList() []models.Webmaster {
	return s.siteDao.GetWebmasterInfoList()
}

func (s *SiteService) CreateWebmaster(webmaster *models.Webmaster) error {
	return s.siteDao.CreateWebmaster(webmaster)
}

func (s *SiteService) UpdateWebmaster(id int64, fields []string, webmaster *models.Webmaster) error {
	return s.siteDao.UpdateWebmaster(id, fields, webmaster)
}
