package admin

import (
	"flygoose/pkg/beans"
	"flygoose/pkg/models"
	"flygoose/pkg/tools"
	"flygoose/web/controllers/comm"
	"flygoose/web/middlers"
	"flygoose/web/services"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"time"
	"unicode/utf8"
)

type BlogController struct {
	comm.BaseComponent
	BlogSrv *services.BlogService
}

func NewBlogController() *BlogController {
	return &BlogController{BlogSrv: services.NewBlogService()}
}

func (c *BlogController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/createBlog", "CreateBlog", middlers.CheckAdminToken)                   //创建博客
	b.Handle("POST", "/updateBlog", "UpdateBlog", middlers.CheckAdminToken)                   //更新博客
	b.Handle("POST", "/updateBlogStatus", "UpdateBlogStatus", middlers.CheckAdminToken)       //更新博客状态
	b.Handle("POST", "/publishBlog", "PublishBlog", middlers.CheckAdminToken)                 //发布博客
	b.Handle("POST", "/getBlogListByStatus", "GetBlogListByStatus", middlers.CheckAdminToken) //根据博客状态获取博客列表
	b.Handle("POST", "/getBlogListByWord", "GetBlogListByWord", middlers.CheckAdminToken)     //根据关键字搜索获取博客列表
	b.Handle("POST", "/searchBlog", "SearchBlog", middlers.CheckAdminToken)                   //搜索博客
	b.Handle("POST", "/getBlogDetail", "GetBlogDetail", middlers.CheckAdminToken)             //獲取博客詳情
}

func (c *BlogController) CreateBlog() {
	var bean beans.CreateBlogBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	blog := models.Blog{
		ID:          tools.GenBlogId(),
		Title:       bean.Title,
		Content:     bean.Content,
		Html:        bean.Html,
		IsHtml:      bean.IsHtml,
		IsTop:       bean.IsTop,
		ReadCount:   0,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		PublishTime: time.Now(),
		Thumbnail:   bean.Thumbnail,
		Status:      models.BlogStatusCreated,
		CateId:      bean.CateId,
		Tags:        bean.Tags,
	}

	if bean.Intro == "" && bean.Content != "" {
		if utf8.RuneCountInString(bean.Content) > 127 {
			bean.Intro = string([]rune(bean.Content)[:126])
		} else {
			bean.Intro = bean.Content
		}
	}

	err := c.BlogSrv.CreateBlog(&blog)
	if err != nil {
		c.RespFailedMessage("创建失败")
	} else {
		c.RespSuccess(blog, "创建成功")
	}
}

func (c *BlogController) UpdateBlog() {
	fields := make([]string, 0)
	blog := &models.Blog{}

	id, exist := c.PostFormInt64Default("id")
	if !exist || id <= 0 {
		c.RespFailedMessage("id参数不存在或者错误")
		return
	}

	title := c.PostFormString("title")
	if title != "" {
		fields = append(fields, "Title")
		blog.Title = title
	}

	intro, exist := c.PostFormStringDefault("intro")
	if exist {
		fields = append(fields, "Intro")
		blog.Intro = intro
	}

	content, exist := c.PostFormStringDefault("content")
	if exist {
		fields = append(fields, "Content")
		blog.Content = content
	}

	html, exist := c.PostFormStringDefault("html")
	if exist {
		fields = append(fields, "Html")
		blog.Html = html
	}

	isHtml, exist := c.PostFormIntDefault("isHtml")
	if exist && (isHtml == 0 || isHtml == 1) {
		fields = append(fields, "IsHtml")
		blog.IsHtml = isHtml
	}

	isTop, exist := c.PostFormIntDefault("isTop")
	if exist && (isTop == 0 || isTop == 1) {
		fields = append(fields, "IsTop")
		blog.IsTop = isTop
	}

	thumbnail, exist := c.PostFormStringDefault("thumbnail")
	if exist {
		fields = append(fields, "Thumbnail")
		blog.Thumbnail = thumbnail
	}

	cateId, exist := c.PostFormInt64Default("cateId")
	if exist {
		fields = append(fields, "CateId")
		blog.CateId = cateId
	}

	tags, exist := c.PostFormStringDefault("tags")
	if exist {
		fields = append(fields, "Tags")
		blog.Tags = tags
	}

	if len(fields) == 0 {
		c.RespFailedMessage("请输入要修改的参数")
		return
	}

	fields = append(fields, "UpdateTime")
	blog.UpdateTime = time.Now()

	err := c.BlogSrv.UpdateBlog(id, fields, blog)
	if err != nil {
		c.RespFailedMessage("更新失败:" + err.Error())
	} else {
		c.RespSuccessMessage("更新成功")
	}
}

func (c *BlogController) UpdateBlogStatus() {
	var bean beans.UpdateBlogStatusBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	err := c.BlogSrv.UpdateBlogStatus(bean.Id, models.BlogStatus(bean.Status))
	if err != nil {
		c.RespFailedMessage("更新博客状态出错")
	} else {
		c.RespSuccessMessage("更新成功")
	}
}

func (c *BlogController) PublishBlog() {
	fields := make([]string, 0)
	blog := &models.Blog{}

	id, exist := c.PostFormInt64Default("id")

	title := c.PostFormString("title")
	if title != "" {
		fields = append(fields, "Title")
		blog.Title = title
	}

	intro, exist := c.PostFormStringDefault("intro")
	if exist {
		fields = append(fields, "Intro")
		blog.Intro = intro
	}

	content, exist := c.PostFormStringDefault("content")
	if exist {
		fields = append(fields, "Content")
		blog.Content = content
	}

	html, exist := c.PostFormStringDefault("html")
	if exist {
		fields = append(fields, "Html")
		blog.Html = html
	}

	isHtml, exist := c.PostFormIntDefault("isHtml")
	if exist && (isHtml == 0 || isHtml == 1) {
		fields = append(fields, "IsHtml")
		blog.IsHtml = isHtml
	}

	isTop, exist := c.PostFormIntDefault("isTop")
	if exist && (isTop == 0 || isTop == 1) {
		fields = append(fields, "IsTop")
		blog.IsTop = isTop
	}

	thumbnail, exist := c.PostFormStringDefault("thumbnail")
	if exist {
		fields = append(fields, "Thumbnail")
		blog.Thumbnail = thumbnail
	}

	cateId, exist := c.PostFormInt64Default("cateId")
	if exist {
		fields = append(fields, "CateId")
		blog.CateId = cateId
	}

	tags, exist := c.PostFormStringDefault("tags")
	if exist {
		fields = append(fields, "Tags")
		blog.Tags = tags
	}

	if len(fields) == 0 {
		c.RespFailedMessage("请输入要修改的参数")
		return
	}

	fields = append(fields, "Status")
	blog.Status = models.BlogStatusCreated

	blog.UpdateTime = time.Now()
	fields = append(fields, "UpdateTime")

	blog.PublishTime = time.Now()
	fields = append(fields, "PublishTime")

	var err error
	if id <= 0 { //没有blog id ,所以新建一个博客,并且把状态改为 models.BlogStatusPublished
		blog.CreateTime = time.Now()
		fields = append(fields, "CreateTime")
		err = c.BlogSrv.CreateBlog(blog)
	} else {
		err = c.BlogSrv.UpdateBlog(id, fields, blog)
	}

	if err != nil {
		c.RespFailedMessage("发布失败")
	} else {
		c.RespSuccessMessage("发布成功")
	}
}

func (c *BlogController) GetBlogListByStatus() {
	var bean beans.BlogPageStatusBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	list, count := c.BlogSrv.GetBlogListByStatus(bean.Status, bean.PageNum, bean.PageSize)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取数据成功")
}

func (c *BlogController) GetBlogListByWord() {
	var bean beans.BlogPageWord
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	list, count := c.BlogSrv.GetBlogListByWord(bean.Word, bean.PageNum, bean.PageSize)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取数据成功")
}

func (c *BlogController) SearchBlog() {
	var bean beans.SearchBlogBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	list, count := c.BlogSrv.SearchBlog(&bean)
	c.RespSuccess(iris.Map{
		"list":  list,
		"count": count,
	}, "获取数据成功")
}

func (c *BlogController) GetBlogDetail() {
	var bean beans.GetBlogDetailBean
	if err := c.Ctx.ReadForm(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	blog, _ := c.BlogSrv.GetBlogDetail(bean.Id)
	if blog == nil {
		c.RespFailedMessage("博客不存在")
	} else {
		c.RespSuccess(blog, "获取博客成功")
	}
}
