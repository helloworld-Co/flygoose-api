package daos

import (
	"flygoose/datasource"
	"flygoose/pkg/models"
	"flygoose/pkg/tlog"
	"gorm.io/gorm"
)

type LinkDao struct {
	db *gorm.DB
}

func NewLinkDao() *LinkDao {
	return &LinkDao{db: datasource.GetMasterDB()}
}

func (dao *LinkDao) Create(link *models.Link) error {
	result := dao.db.Create(link)
	if result.Error != nil {
		tlog.Error2("LinkDao:Create 出错", result.Error)
	}
	return result.Error
}

func (dao *LinkDao) Update(id int64, fields []string, link *models.Link) error {
	result := dao.db.Model(&models.Link{}).Select(fields).Where("id=?", id).Updates(link)
	if result.Error != nil {
		tlog.Error2("LinkDao:Update 出错", result.Error)
	}
	return result.Error
}

func (dao *LinkDao) GetAllLinkList(num int, size int) ([]models.Link, int64) {
	var list []models.Link
	var count int64

	result := dao.db.Order("seq ").Order("status desc").Limit(size).Offset((num - 1) * size).Find(&list)
	if result.Error != nil {
		tlog.Error2("LinkDao:GetAllLinkList Find 出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Link{}).Count(&count)
	if result.Error != nil {
		tlog.Error2("LinkDao:GetAllLinkList Count 出错", result.Error)
		return nil, 0
	}

	return list, count
}

func (dao *LinkDao) GetLinkListByStatus(status int, num int, size int) ([]models.Link, int64) {
	var list []models.Link
	var count int64

	result := dao.db.Order("seq ").Order("status desc").Where("status=?", status).Limit(size).Offset((num - 1) * size).Find(&list)
	if result.Error != nil {
		tlog.Error2("LinkDao:GetLinkListByStatus Find 出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Link{}).Where("status=?", status).Count(&count)
	if result.Error != nil {
		tlog.Error2("LinkDao:GetLinkListByStatus Count 出错", result.Error)
		return nil, 0
	}

	return list, count
}

func (dao *LinkDao) GetLinkInfo(id int64) *models.Link {
	var link models.Link
	result := dao.db.Where("id=?", id).First(&link)
	if result.Error != nil {
		tlog.Error2("LinkDao:GetLinkInfo 出错", result.Error)
		return nil
	}
	return &link
}

func (dao *LinkDao) GetTotal() int64 {
	var total int64
	result := dao.db.Model(&models.Link{}).Where("status = ? ", models.LinkStatusNormal).Count(&total)
	if result.Error != nil {
		tlog.Error2("LinkDao GetTotal出错", result.Error)
		return 0
	}
	return total
}
