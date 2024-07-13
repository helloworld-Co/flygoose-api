package services

import (
	"flygoose/pkg/models"
	"flygoose/web/daos"
)

type BannerService struct {
	bannerDao *daos.BannerDao
}

func NewBannerService() *BannerService {
	return &BannerService{bannerDao: daos.NewBannerDao()}
}

func (srv *BannerService) Create(m *models.Banner) error {
	return srv.bannerDao.Create(m)
}

func (srv *BannerService) Update(id int64, fields []string, banner *models.Banner) error {
	return srv.bannerDao.Update(id, fields, banner)
}

func (srv *BannerService) GetBannerInfo(id int64) *models.Banner {
	return srv.bannerDao.GetBannerInfo(id)
}

func (srv *BannerService) GetBannerList(status int, num int, size int) ([]models.Banner, int64) {
	if status == -1 { //全部
		return srv.bannerDao.GetAllBannerList(num, size)
	} else {
		return srv.bannerDao.GetBannerListByStatus(status, num, size)
	}
}
