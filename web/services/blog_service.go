package services

import (
	"errors"
	"flygoose/pkg/beans"
	"flygoose/pkg/models"
	"flygoose/pkg/tools"
	"flygoose/web/daos"
	"strings"
	"time"
)

type BlogService struct {
	blogDao *daos.BlogDao
}

func NewBlogService() *BlogService {
	return &BlogService{blogDao: daos.NewBlogDao()}
}

func (s *BlogService) GetBlogListByStatus(status int, pageNum int, pageSize int) ([]models.Blog, int64) {
	var list []models.Blog
	var count int64
	if status == 0 {
		list, count = s.blogDao.GetAllBlogList(pageNum, pageSize)
	} else {
		list, count = s.blogDao.GetBlogListByStatus(models.BlogStatus(status), pageNum, pageSize)
	}

	return list, count
}

func (s *BlogService) CreateBlog(blog *models.Blog) error {
	return s.blogDao.CreateBlog(blog)
}

func (s *BlogService) UpdateBlog(id int64, fields []string, blog *models.Blog) error {
	_, exist := tools.Find(fields, "CateId")
	if !exist {
		return s.blogDao.UpdateBlog(id, fields, blog)
	}

	cateDao := daos.NewCategoryDao()
	cate := cateDao.GetById(blog.CateId)
	if cate == nil {
		if blog.CateId > 0 {
			return errors.New("CateId对应的分类不存在")
		} else {
			blog.CateId = 0
		}
	} else {
		if cate.Status == models.CateStatusOffline {
			return errors.New("CateId对应的分类已下架")
		}
	}

	return s.blogDao.UpdateBlog(id, fields, blog)
}

func (s *BlogService) UpdateBlogStatus(blogId int64, status models.BlogStatus) error {
	fields := make([]string, 0)
	blog := models.Blog{}

	fields = append(fields, "Status")
	blog.Status = status

	fields = append(fields, "UpdateTime")
	blog.UpdateTime = time.Now()

	if status == models.BlogStatusPublished {
		blog.UpdateTime = time.Now()
		fields = append(fields, "UpdateTime")
	}

	return s.blogDao.UpdateBlog(blogId, fields, &blog)
}

func (s *BlogService) GetBlogListByWord(word string, pageNum int, pageSize int) ([]models.Blog, int64) {
	return s.blogDao.GetBlogListByWord(word, pageNum, pageSize)
}

func (s *BlogService) SearchBlog(b *beans.SearchBlogBean) ([]models.Blog, int64) {
	return s.blogDao.SearchBlog(b.Word, b.Status, b.PageNum, b.PageSize)
}

func (s *BlogService) GetBlogListByAction(action int, num int, size int) ([]models.Blog, bool) {
	return s.blogDao.GetBlogListByAction(action, num, size)
}

func (s *BlogService) GetCateBlogList(cateId int64, num int, size int) ([]models.Blog, bool) {
	return s.blogDao.GetCateBlogList(cateId, num, size)
}

func (s *BlogService) GetAllTags() []string {
	tags := s.blogDao.GetAllTags()

	sets := make(map[string]bool)

	for _, r := range tags {
		if r != "" {
			arr := strings.Split(r, ",")
			if len(arr) > 0 {
				for _, t := range arr {
					if t != "" {
						sets[t] = true
					}
				}
			}
		}
	}

	var list []string
	for tag := range sets { // Loop
		list = append(list, tag)
	}

	return list
}

func (s *BlogService) GetBlogListByTag(tagName string, num int, size int) ([]models.Blog, bool) {
	return s.blogDao.GetBlogListByTag(tagName, num, size)
}

func (s *BlogService) GetBlogDetail(blogId int64) (*models.Blog, []models.Blog) {
	return s.blogDao.GetBlogDetail(blogId), s.GetAuthorOtherBlog(blogId)
}

// 获取作者的其它博客
func (s *BlogService) GetAuthorOtherBlog(blogId int64) []models.Blog {
	return s.blogDao.GetAuthorOtherBlog(blogId)
}
