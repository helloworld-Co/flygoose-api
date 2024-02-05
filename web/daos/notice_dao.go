package daos

import (
	"flygoose/datasource"
	"flygoose/pkg/models"
	"flygoose/pkg/tlog"
	"gorm.io/gorm"
	"time"
)

type NoticeDao struct {
	db *gorm.DB
}

func NewNoticeDao() *NoticeDao {
	return &NoticeDao{db: datasource.GetMasterDB()}
}

func (dao *NoticeDao) Create(m *models.Notice) error {
	result := dao.db.Create(m)
	if result.Error != nil {
		tlog.Error2("NoticeDao:Create 出错", result.Error)
	} else {
		//mengxianhou@20230124 公告只能有一个上架的 其他的全部更新为下架
		tx := dao.db.Model(&models.Notice{}).Exec("update notice set status = 0,update_time = ? where id != ? and status = 1 ", time.Now(), m.ID)
		if tx.Error != nil {
			tlog.Error2("NoticeDao:更新状态 出错", tx.Error)
			return tx.Error
		}
	}
	return result.Error
}

func (dao *NoticeDao) Update(id int64, fields []string, m *models.Notice) error {
	result := dao.db.Model(&models.Notice{}).Where("id=?", id).Select(fields).Updates(m)
	if result.Error != nil {
		tlog.Error2("NoticeDao:Update 出错", result.Error)
	} else {
		if m.Status == 1 {
			//mengxianhou@20230124 公告只能有一个上架的 其他的全部更新为下架
			dao.db.Model(&models.Notice{}).Exec("update notice set status = 0,update_time = ? where id != ?", time.Now(), id)
		}
	}
	return result.Error
}

func (dao *NoticeDao) GetAllNoticeList(num int, size int) ([]models.Notice, int64) {
	var list []models.Notice
	var count int64

	result := dao.db.Limit(size).Order("create_time desc").Offset(size * (num - 1)).Find(&list)
	if result.Error != nil {
		tlog.Error2("NoticeDao:GetAllNoticeList Find 出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Notice{}).Count(&count)
	if result.Error != nil {
		tlog.Error2("NoticeDao:GetAllNoticeList Count 出错", result.Error)
		return nil, 0
	}

	return list, count
}

func (dao *NoticeDao) GetNoticeListByStatus(status int, num int, size int) ([]models.Notice, int64) {
	var list []models.Notice
	var count int64

	result := dao.db.Where("status=?", status).Order("create_time desc").Limit(size).Offset(size * (num - 1)).Find(&list)
	if result.Error != nil {
		tlog.Error2("NoticeDao:GetNoticeListByStatus Find 出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Notice{}).Where("status=?", status).Count(&count)
	if result.Error != nil {
		tlog.Error2("NoticeDao:GetNoticeListByStatus Count 出错", result.Error)
		return nil, 0
	}

	return list, count
}

func (dao *NoticeDao) GetNoticeInfo(id int64) *models.Notice {
	var notice models.Notice
	result := dao.db.Where("id=?", id).First(&notice)
	if result.Error != nil {
		tlog.Error2("NoticeDao:GetNoticeInfo Count 出错", result.Error)
		return nil
	}
	return &notice
}
