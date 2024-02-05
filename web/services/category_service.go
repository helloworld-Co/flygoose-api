package services

import (
	"flygoose/pkg/models"
	"flygoose/web/daos"
)

type CategoryService struct {
	categoryDao *daos.CategoryDao
}

func NewCategoryService() *CategoryService {
	return &CategoryService{categoryDao: daos.NewCategoryDao()}
}

func (s *CategoryService) Create(m *models.Category) error {
	return s.categoryDao.Create(m)
}

func (s *CategoryService) Update(id int64, fields []string, m *models.Category) error {
	return s.categoryDao.Update(id, fields, m)
}

func (s *CategoryService) GetCategoryList(status int) []models.Category {
	if status == -1 {
		return s.categoryDao.GetAllCategoryList()
	} else {
		return s.categoryDao.GetCategoryListByStatus(status)
	}
}
