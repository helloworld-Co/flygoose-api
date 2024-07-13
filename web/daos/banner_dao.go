package daos

import (
	"flygoose/datasource"
	"flygoose/pkg/models"
	"flygoose/pkg/tlog"
	"gorm.io/gorm"
)

type BannerDao struct {
	db *gorm.DB
}

func NewBannerDao() *BannerDao {
	return &BannerDao{db: datasource.GetMasterDB()}
}

func (dao *BannerDao) Create(m *models.Banner) error {
	result := dao.db.Create(m)
	if result.Error != nil {
		tlog.Error2("BannerDao:Create出错", result.Error)
	}
	return result.Error
}

func (dao *BannerDao) Update(id int64, fields []string, banner *models.Banner) error {
	result := dao.db.Model(&models.Banner{}).Select(fields).Where("id=?", id).Updates(banner)
	if result.Error != nil {
		tlog.Error2("BannerDao:Update出错", result.Error)
	}
	return result.Error
}

func (dao *BannerDao) GetBannerInfo(id int64) *models.Banner {
	var banner models.Banner
	result := dao.db.Where("id=?", id).First(&banner)
	if result.Error != nil {
		tlog.Error2("BannerDao:GetBannerInfo出错", result.Error)
		return nil
	}
	return &banner
}

func (dao *BannerDao) GetAllBannerList(num int, size int) ([]models.Banner, int64) {
	var list []models.Banner
	var count int64

	result := dao.db.Order("seq ").Order("status desc").Limit(size).Offset((num - 1) * size).Find(&list)
	if result.Error != nil {
		tlog.Error2("BannerDao:GetAllBannerList Find操作出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Banner{}).Count(&count)
	if result.Error != nil {
		tlog.Error2("BannerDao:GetAllBannerList Count操作出错", result.Error)
		return nil, 0
	}

	return list, count
}

func (dao *BannerDao) GetBannerListByStatus(status int, num int, size int) ([]models.Banner, int64) {
	var list []models.Banner
	var count int64

	result := dao.db.Order("seq ").Order("status desc").Where("status=?", status).Limit(size).Offset((num - 1) * size).Find(&list)
	if result.Error != nil {
		tlog.Error2("BannerDao:GetBannerListByStatus Find操作出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Banner{}).Where("status=?", status).Count(&count)
	if result.Error != nil {
		tlog.Error2("BannerDao:GetBannerListByStatus Count操作出错", result.Error)
		return nil, 0
	}

	return list, count
}

func (dao *BannerDao) GetTotal() int64 {
	var total int64
	result := dao.db.Model(&models.Banner{}).Where("status = ?", models.LinkStatusNormal).Count(&total)
	if result.Error != nil {
		tlog.Error2("BannerDao GetTotal报错", result.Error)
		return 0
	}

	return total
}
