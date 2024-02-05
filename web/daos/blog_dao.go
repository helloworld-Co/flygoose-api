package daos

import (
	"flygoose/datasource"
	"flygoose/pkg/models"
	"flygoose/pkg/tlog"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type BlogDao struct {
	db *gorm.DB
}

func NewBlogDao() *BlogDao {
	return &BlogDao{db: datasource.GetMasterDB()}
}

// 获取所有状态的博客
func (dao *BlogDao) GetAllBlogList(pageNum int, pageSize int) ([]models.Blog, int64) {
	var list []models.Blog
	var count int64

	result := dao.db.Omit("Content", "Html").Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("create_time desc").Find(&list)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetAllBlogList Find 出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Blog{}).Count(&count)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetAllBlogList Count 出错", result.Error)
		return nil, 0
	}

	return list, count
}

// 获取某一种状态的博客的列表
func (dao *BlogDao) GetBlogListByStatus(status models.BlogStatus, pageNum int, pageSize int) ([]models.Blog, int64) {
	var list []models.Blog
	var count int64

	result := dao.db.Omit("Content", "Html").Where("status=?", status).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("create_time desc").Find(&list)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetBlogListByStatus Find 出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Blog{}).Where("status=?", status).Count(&count)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetBlogListByStatus Count 出错", result.Error)
		return nil, 0
	}

	return list, count
}

func (dao *BlogDao) CreateBlog(blog *models.Blog) error {
	result := dao.db.Create(&blog)
	if result.Error != nil {
		tlog.Error2("BlogDao:CreateBlog 出错", result.Error)
	}
	return result.Error
}

func (dao *BlogDao) UpdateBlog(id int64, fields []string, blog *models.Blog) error {
	result := dao.db.Model(&models.Blog{}).Select(fields).Where("id=?", id).Updates(blog)
	if result.Error != nil {
		tlog.Error2("BlogDao:UpdateBlog 出错", result.Error)
	}
	return result.Error
}

func (dao *BlogDao) GetBlogListByWord(word string, pageNum int, pageSize int) ([]models.Blog, int64) {
	var list []models.Blog
	var count int64

	sqlCount := `select count(*) as count from blog where title like '%` + word + `%'`
	result := dao.db.Raw(sqlCount).Scan(&count)
	if result.Error != nil || count == 0 {
		tlog.Error2("BlogDao:GetBlogListByWord Count 出错", result.Error)
		return nil, 0
	}

	sql := `select id,title,intro,is_top,read_count,create_time,update_time,publish_time,thumbnail,status,cate_id,tags 
from blog where title like '%` + word + `%' order by create_time desc limit ? offset ?`
	result = dao.db.Model(&models.Blog{}).Raw(sql, pageSize, pageSize*(pageNum-1)).Scan(&list)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetBlogListByWord Scan 出错", result.Error)
		return nil, 0
	}

	return list, 0
}

func (dao *BlogDao) SearchBlog(word string, status int, num int, size int) ([]models.Blog, int64) {
	var sql = ""
	var sqlCount = ""
	if word != "" && status > 0 {
		sql = "select * from blog where title like " + "'%" + word + "%'" + " and status=" + strconv.Itoa(status) +
			" order by create_time desc  limit " + strconv.Itoa(size*(num-1)) + " , " + strconv.Itoa(size)
		sqlCount = "select count(*) as count from blog where title like " + "'%" + word + "%'" + " and status=" + strconv.Itoa(status)
	} else if word != "" && status <= 0 {
		sql = "select * from blog where title like " + "'%" + word + "%'" + " order by create_time desc limit " + strconv.Itoa(size*(num-1)) + " , " + strconv.Itoa(size)
		sqlCount = "select count(*) as count from blog where title like " + "'%" + word + "%'"
	} else if status > 0 && word == "" {
		s1 := "select * from blog where  status=%d order by create_time desc  limit %d , %d"
		sql = fmt.Sprintf(s1, status, size*(num-1), size)

		s2 := "select count(*) as count from blog where status=%d"
		sqlCount = fmt.Sprintf(s2, status)
	} else { //都不传，获取全部
		sql = "select * from blog order by create_time desc  limit  " + strconv.Itoa(size*(num-1)) + " , " + strconv.Itoa(size)
		sqlCount = "select count(*) as count from blog"
	}

	var list []models.Blog
	var count int64

	result := dao.db.Raw(sqlCount).Scan(&count)
	if result.Error != nil || count == 0 {
		tlog.Error2("BlogDao:SearchBlog Count 出错", result.Error)
		return nil, 0
	}

	result = dao.db.Model(&models.Blog{}).Raw(sql).Scan(&list)
	if result.Error != nil {
		tlog.Error2("BlogDao:SearchBlog Scan 出错", result.Error)
		return nil, 0
	}

	return list, count
}

func (dao *BlogDao) GetBlogListByAction(action int, num int, size int) ([]models.Blog, bool) {
	var sql string
	if action == 0 { //最新发布
		sql = fmt.Sprintf("select id,title,intro,thumbnail,publish_time,read_count,is_top,tags from blog where status=%d order by is_top desc ,  publish_time desc limit %d,%d", models.BlogStatusPublished, (num-1)*size, size)
	} else { //阅读数最多
		sql = fmt.Sprintf("select id,title,intro,thumbnail,publish_time,read_count,is_top,tags from blog where status=%d order by is_top desc ,  read_count desc , publish_time desc limit %d,%d", models.BlogStatusPublished, (num-1)*size, size)
	}

	var list []models.Blog
	var hasMore bool

	result := dao.db.Model(&models.Blog{}).Raw(sql).Scan(&list)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetBlogListByAction Scan 出错", result.Error)
		return nil, false
	}

	hasMore = size == len(list)
	return list, hasMore
}

func (dao *BlogDao) GetCateBlogList(cateId int64, num int, size int) ([]models.Blog, bool) {
	var list []models.Blog
	var hasMore bool

	sql := fmt.Sprintf(" select id,title,intro,thumbnail,publish_time,read_count,is_top,tags from blog where status=%d and cate_id=%d order by is_top desc , publish_time desc limit %d,%d", models.BlogStatusPublished, cateId, (num-1)*size, size)
	result := dao.db.Model(&models.Blog{}).Raw(sql).Scan(&list)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetCateBlogList Scan 出错", result.Error)
		return nil, false
	}

	hasMore = size == len(list)
	return list, hasMore
}

func (dao *BlogDao) GetAllTags() []string {
	var tags []string
	sql := fmt.Sprintf("select tags from blog where status=%d", models.BlogStatusPublished)
	result := dao.db.Model(&models.Blog{}).Raw(sql).Scan(&tags)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetAllTags Scan 出错", result.Error)
		return nil
	}

	return tags
}

func (dao *BlogDao) GetBlogListByTag(tagName string, num int, size int) ([]models.Blog, bool) {
	sql := "select * from blog where tags like " + "'%" + tagName + "%'" + " and status=" + strconv.Itoa(int(models.BlogStatusPublished)) +
		" order by publish_time desc limit " + strconv.Itoa(size*(num-1)) + " , " + strconv.Itoa(size)

	var list []models.Blog
	var hasMore bool
	result := dao.db.Model(&models.Blog{}).Raw(sql).Scan(&list)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetBlogListByTag Scan 出错", result.Error)
		return nil, false
	}

	hasMore = size == len(list)
	return list, hasMore
}

func (dao *BlogDao) GetBlogDetail(blogId int64) *models.Blog {
	var blog models.Blog
	result := dao.db.First(&blog, blogId)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetBlogDetail Scan 出错", result.Error)
		return nil
	}
	return &blog
}

func (dao *BlogDao) GetAuthorOtherBlog(id int64) []models.Blog {
	var list []models.Blog
	result := dao.db.Model(&models.Blog{}).Where("id != ? and status=?", id, models.BlogStatusPublished).Order("read_count desc,create_time desc").Limit(10).Find(&list)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetAuthorOtherBlog Scan 出错", result.Error)
		return nil
	}
	return list
}

func (dao *BlogDao) GetTotal() int64 {
	var total int64
	result := dao.db.Model(&models.Blog{}).Where("status = ?", models.BlogStatusPublished).Count(&total)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetTotal count 出错", result.Error)
		return 0
	}
	return total
}

func (dao *BlogDao) GetTodayTotal() int64 {
	var total int64
	now := time.Now()
	todayStart := now.Format("2006-01-02") + " 00:00:00"
	todayEnd := now.Format("2006-01-02") + " 23:59:59"
	result := dao.db.Model(&models.Blog{}).Where("status = ? and create_time >= ? and create_time <= ? ", models.BlogStatusPublished, todayStart, todayEnd).Count(&total)
	if result.Error != nil {
		tlog.Error2("BlogDao:GetTotal count 出错", result.Error)
		return 0
	}
	return total
}
