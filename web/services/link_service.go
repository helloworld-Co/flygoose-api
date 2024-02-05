package services

import (
	"flygoose/pkg/models"
	"flygoose/web/daos"
)

type LinkService struct {
	linkDao *daos.LinkDao
}

func NewLinkService() *LinkService {
	return &LinkService{linkDao: daos.NewLinkDao()}
}

func (s *LinkService) Create(m *models.Link) error {
	return s.linkDao.Create(m)
}

func (s *LinkService) Update(id int64, fields []string, link *models.Link) error {
	return s.linkDao.Update(id, fields, link)
}

func (s *LinkService) GetLinkList(status int, num int, size int) ([]models.Link, int64) {
	if status == -1 { //全部
		return s.linkDao.GetAllLinkList(num, size)
	} else {
		return s.linkDao.GetLinkListByStatus(status, num, size)
	}
}

func (s *LinkService) GetLinkInfo(id int64) *models.Link {
	return s.linkDao.GetLinkInfo(id)
}
