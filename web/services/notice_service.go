package services

import (
	"flygoose/pkg/models"
	"flygoose/web/daos"
)

type NoticeService struct {
	noticeDao *daos.NoticeDao
}

func NewNoticeService() *NoticeService {
	return &NoticeService{noticeDao: daos.NewNoticeDao()}
}

func (s *NoticeService) Create(m *models.Notice) error {
	return s.noticeDao.Create(m)
}

func (s *NoticeService) Update(id int64, fields []string, m *models.Notice) error {
	return s.noticeDao.Update(id, fields, m)
}

func (s *NoticeService) GetNoticeList(status int, num int, size int) ([]models.Notice, int64) {
	if status == -1 {
		return s.noticeDao.GetAllNoticeList(num, size)
	} else {
		return s.noticeDao.GetNoticeListByStatus(status, num, size)
	}
}

func (s *NoticeService) GetNoticeInfo(id int64) *models.Notice {
	return s.noticeDao.GetNoticeInfo(id)

}
