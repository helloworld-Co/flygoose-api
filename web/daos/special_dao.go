package daos

import (
	"errors"
	"flygoose/datasource"
	"flygoose/pkg/models"
	"flygoose/pkg/tlog"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type SpecialDao struct {
	db *gorm.DB
}

func NewSpecialDao() *SpecialDao {
	return &SpecialDao{db: datasource.GetMasterDB()}
}

func (dao *SpecialDao) Create(m *models.Special) error {
	result := dao.db.Create(m)
	if result.Error != nil {
		tlog.Error2("SpecialDao:Create 出错", result.Error)
	}
	return result.Error
}

func (dao *SpecialDao) First(id int64) (*models.Special, error) {
	var special models.Special
	result := dao.db.First(&special, "id=?", id)
	if result.Error != nil {
		tlog.Error2("SpecialDao:First 出错", result.Error)
		return nil, result.Error
	}

	return &special, nil
}

func (dao *SpecialDao) Update(id int64, fields []string, special *models.Special) error {
	result := dao.db.Model(&models.Special{}).Select(fields).Where("id=?", id).Updates(special)
	if result.Error != nil {
		tlog.Error2("SpecialDao:Update 出错", result.Error)
	}
	return result.Error
}

func (dao *SpecialDao) AddSection(section *models.Section) error {
	special, err := dao.First(section.SpecialId)
	if err != nil || special == nil || special.Status == models.SpecialStatusDeleted {
		tlog.Error2("SpecialDao:AddSection First 出错", err)
		return errors.New("专栏不存在")
	}

	result := dao.db.Create(section)
	if result.Error != nil {
		tlog.Error2("SpecialDao:AddSection Create 出错", result.Error)
	}
	return result.Error
}

func (dao *SpecialDao) UpdateSection(specialId int64, sectionId int64, fields []string, section *models.Section) error {
	result := dao.db.Model(&models.Section{}).Select(fields).Where("id=? and special_id=?", sectionId, specialId).Updates(section)
	if result.Error != nil {
		tlog.Error2("SpecialDao:UpdateSection 出错", result.Error)
	}
	return result.Error
}

func (dao *SpecialDao) SearchSpecial(word string, status int, num int, size int) ([]models.Special, int64) {
	var sql = ""
	var sqlCount = ""
	if word != "" && status > 0 {
		sql = "select * from special where title like " + "'%" + word + "%'" + " and status=" + strconv.Itoa(status) +
			" order by create_time desc  limit " + strconv.Itoa(size*(num-1)) + " , " + strconv.Itoa(size)
		sqlCount = "select count(*) as count from special where title like " + "'%" + word + "%'" + " and status=" + strconv.Itoa(status)
	} else if word != "" && status <= 0 {
		sql = "select * from special where title like " + "'%" + word + "%'" + " order by create_time desc limit " + strconv.Itoa(size*(num-1)) + " , " + strconv.Itoa(size)
		sqlCount = "select count(*) as count from special where title like " + "'%" + word + "%'"
	} else if status > 0 && word == "" {
		s1 := "select * from special where  status=%d order by create_time desc  limit %d , %d"
		sql = fmt.Sprintf(s1, status, size*(num-1), size)

		s2 := "select count(*) as count from special where status=%d"
		sqlCount = fmt.Sprintf(s2, status)
	} else { //都不传，获取全部
		sql = "select * from special order by create_time desc  limit  " + strconv.Itoa(size*(num-1)) + " , " + strconv.Itoa(size)
		sqlCount = "select count(*) as count from special"
	}

	var list []models.Special
	var count int64

	result := dao.db.Raw(sqlCount).Scan(&count)
	if result.Error != nil || count == 0 {
		tlog.Error2("SpecialDao:SearchSpecial Count 出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Special{}).Raw(sql).Scan(&list)
	if result.Error != nil {
		tlog.Error2("SpecialDao:SearchSpecial List 出错", result.Error)
		return nil, 0
	}

	return list, count
}

// 获取专栏的已发布的小节数以及专栏的总的阅读数
func (dao *SpecialDao) GetSpecialPublishedCountAndReadCount(specialId int64) (int64, int64) {
	var sqlReadCount = fmt.Sprintf("select read_count as count from section where special_id=%d  and status=%d", specialId, models.SectionStatusPublished)
	var sqlPublishedCount = fmt.Sprintf("select id from section where status=%d and special_id=%d", models.SectionStatusPublished, specialId)

	var readCount int64 = 0
	var publishedCount int64 = 0

	var list []string
	result := dao.db.Raw(sqlReadCount).Find(&list)
	if result.Error != nil {
		tlog.Error2("SpecialDao:GetSpecialPublishedCountAndReadCount readCount 出错", result.Error)
		return 0, 0
	}

	for _, r := range list {
		count, _ := strconv.ParseInt(r, 10, 64)
		readCount = readCount + count
	}

	var list2 []string
	result = dao.db.Raw(sqlPublishedCount).Find(&list2)
	if result.Error != nil {
		tlog.Error2("SpecialDao:GetSpecialPublishedCountAndReadCount publishedCount 出错", result.Error)
		return 0, 0
	}

	publishedCount = int64(len(list2))

	return publishedCount, readCount
}

func (dao *SpecialDao) GetSectionTotalAndPublishedCount(specialId int64) (int, int) {
	var sqlTotal = fmt.Sprintf("select count(*) as count from section where special_id=%d", specialId)
	var sqlPublished = fmt.Sprintf("select count(*) as count from section where status=%d and special_id=%d", models.SectionStatusPublished, specialId)

	var totalCount = 0
	var publishedCount = 0
	result := dao.db.Raw(sqlTotal).Scan(&totalCount)
	if result.Error != nil {
		tlog.Error2("SpecialDao:GetSectionTotalAndPublishedCount totalCount 出错", result.Error)
	}

	result = dao.db.Raw(sqlPublished).Scan(&publishedCount)
	if result.Error != nil {
		tlog.Error2("SpecialDao:GetSectionTotalAndPublishedCount publishedCount 出错", result.Error)
	}
	return totalCount, publishedCount
}

func (dao *SpecialDao) GetSectionList(specialId int64, status int) ([]models.Section, int64) {
	var sql, sqlCount = "", ""
	var list, count = make([]models.Section, 0), 0

	if status == 0 {
		sql = fmt.Sprintf("select * from section where special_id=%d order by seq", specialId)
		sqlCount = fmt.Sprintf("select count(*) as count from section where special_id=%d", specialId)
	} else {
		sql = fmt.Sprintf("select * from section where special_id=%d and status=%d order by seq", specialId, status)
		sqlCount = fmt.Sprintf("select count(*) as count from section where special_id=%d and status=%d", specialId, status)
	}

	result := dao.db.Raw(sqlCount).Scan(&count)
	if result.Error != nil || count == 0 {
		tlog.Error2("SpecialDao:GetSectionList Count 出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Section{}).Raw(sql).Scan(&list)
	if result.Error != nil {
		tlog.Error2("SpecialDao:GetSectionList List 出错", result.Error)
		return nil, 0
	}

	return list, int64(count)

}

func (dao *SpecialDao) GetSpecialList(num int, size int) ([]models.Special, bool) {
	var list []models.Special
	var hasMore bool

	sql := fmt.Sprintf("select * from special where status=%d order by publish_time desc limit %d,%d", models.SpecialStatusPublished, size*(num-1), size)
	result := dao.db.Model(&models.Special{}).Raw(sql).Scan(&list)
	if result.Error != nil {
		tlog.Error2("SpecialDao:GetSpecialList List 出错", result.Error)
		return nil, false
	}

	hasMore = size == len(list)
	return list, hasMore
}

func (dao *SpecialDao) FirstBySectionId(sectionId int64) (*models.Section, error) {
	var section models.Section
	result := dao.db.First(&section, sectionId)
	if result.Error != nil {
		tlog.Error2("SpecialDao:FirstBySectionId List 出错", result.Error)
		return nil, result.Error
	}
	return &section, nil
}

func (dao *SpecialDao) GetTotal() int64 {
	var total int64
	result := dao.db.Model(&models.Special{}).Where("status = ?", models.SpecialStatusPublished).Count(&total)
	if result.Error != nil {
		tlog.Error2("SpecialDao GetTotal出错", result.Error)
		return 0
	}
	return total
}

func (dao *SpecialDao) GetTodayTotal() int64 {
	var total int64
	now := time.Now()
	todayStart := now.Format("2006-01-02") + " 00:00:00"
	todayEnd := now.Format("2006-01-02") + " 23:59:59"
	result := dao.db.Model(&models.Special{}).Where("status = ? and create_time >= ? and create_time <= ?", models.SpecialStatusPublished, todayStart, todayEnd).Count(&total)
	if result.Error != nil {
		tlog.Error2("SpecialDao GetTotal出错", result.Error)
		return 0
	}
	return total
}
