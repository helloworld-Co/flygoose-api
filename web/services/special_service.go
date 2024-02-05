package services

import (
	"flygoose/pkg/beans"
	"flygoose/pkg/models"
	"flygoose/pkg/tools"
	"flygoose/web/daos"
	"time"
)

type SpecialService struct {
	specialDao *daos.SpecialDao
}

func NewSpecialService() *SpecialService {
	return &SpecialService{specialDao: daos.NewSpecialDao()}
}

func (s *SpecialService) Create(param *beans.CreateSpecialParam) (*models.Special, error) {
	var special = models.Special{
		ID:          tools.GenSpecialId(),
		Title:       param.Title,
		Intro:       param.Intro,
		Cover:       param.Cover,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		PublishTime: time.Now(),
		Status:      models.SpecialStatusCreated,
	}

	err := s.specialDao.Create(&special)
	if err != nil {
		return nil, err
	} else {
		return &special, nil
	}
}

func (s *SpecialService) Update(specialId int64, fields []string, special *models.Special) error {
	return s.specialDao.Update(specialId, fields, special)
}

func (s *SpecialService) AddSection(param *beans.AddSectionParam) (*models.Section, error) {
	var section = models.Section{
		ID:          tools.GenSectionId(),
		Title:       param.Title,
		SpecialId:   param.SpecialId,
		Seq:         param.Seq,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		PublishTime: time.Now(),
		Status:      models.SectionStatusCreated,
	}

	err := s.specialDao.AddSection(&section)
	if err != nil {
		return nil, err
	} else {
		return &section, nil
	}
}

func (s *SpecialService) UpdateSection(specialId int64, sectionId int64, fields []string, section *models.Section) error {
	return s.specialDao.UpdateSection(specialId, sectionId, fields, section)
}

func (s *SpecialService) SearchSpecial(b *beans.SearchSpecialParam) ([]beans.SpecialBean, int64) {
	list, count := s.specialDao.SearchSpecial(b.Word, b.Status, b.PageNum, b.PageSize)
	if count <= 0 {
		return nil, 0
	}

	sbList := make([]beans.SpecialBean, 0)
	for _, special := range list {
		var sb beans.SpecialBean
		sb.Special = special
		sb.TotalCount, sb.PublishedCount = s.specialDao.GetSectionTotalAndPublishedCount(special.ID)
		sbList = append(sbList, sb)
	}

	return sbList, count
}

func (s *SpecialService) GetSectionList(specialId int64, status int) ([]models.Section, int64) {
	return s.specialDao.GetSectionList(specialId, status)
}

func (s *SpecialService) GetSpecialList(num int, size int) ([]beans.SpecialHomeBean, bool) {
	list, hasMore := s.specialDao.GetSpecialList(num, size)
	if len(list) == 0 {
		return nil, hasMore
	}

	var resultList []beans.SpecialHomeBean
	for _, r := range list {
		var bb beans.SpecialHomeBean
		bb.Special = r
		bb.PublishedCount, bb.ReadCount = s.specialDao.GetSpecialPublishedCountAndReadCount(r.ID)
		resultList = append(resultList, bb)
	}

	return resultList, hasMore
}

func (s *SpecialService) GetSpecialDetail(specialId int64) (*models.Special, []models.Section) {
	special, err := s.specialDao.First(specialId)
	if err != nil || special == nil || special.Status != models.SpecialStatusPublished {
		return nil, nil
	}

	list, _ := s.specialDao.GetSectionList(specialId, int(models.SectionStatusPublished))

	return special, list
}

func (s *SpecialService) GetSectionDetail(sectionId int64) (*models.Section, error) {
	return s.specialDao.FirstBySectionId(sectionId)
}
