package flygoose

import (
	"flygoose/pkg/beans"
	"flygoose/pkg/models"
	"flygoose/web/controllers/comm"
	"flygoose/web/services"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type BlogController struct {
	comm.BaseComponent
}

func NewBlogController() *BlogController {
	return &BlogController{}
}

func (c *BlogController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/getCateList", "GetCateList")                 //获取博客分类列表
	b.Handle("POST", "/getCateBlogList", "GetCateBlogList")         //获取分类下的博客列表
	b.Handle("POST", "/getBlogListByAction", "GetBlogListByAction") //根据action获取博客列表
	b.Handle("POST", "/getAllTags", "GetAllTags")                   //获取所有的标签
	b.Handle("POST", "/getBlogListByTag", "GetBlogListByTag")       //根据标签获取博客
	b.Handle("POST", "/searchBlog", "SearchBlog")                   //搜索博客,只搜索标题
	b.Handle("POST", "/getBlogDetail", "GetBlogDetail")             //获取博客详情
}

func (c *BlogController) SearchBlog() {
	var param beans.SearchBlogBean
	if err := c.Ctx.ReadJSON(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	param.Status = int(models.BlogStatusPublished)

	srv := services.NewBlogService()
	list, _ := srv.SearchBlog(&param)

	hasMore := param.PageSize == len(list)
	c.RespSuccess(iris.Map{
		"list":    list,
		"hasMore": hasMore,
	}, "获取成功")
}

func (c *BlogController) GetBlogListByTag() {
	type Param struct {
		Name     string `json:"name" validate:"required"`     //标签名字
		PageNum  int    `json:"pageNum" validate:"required"`  //第几页，从1开始
		PageSize int    `json:"pageSize" validate:"required"` //每页的数量
	}

	var param Param
	if err := c.Ctx.ReadJSON(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	srv := services.NewBlogService()
	list, hasMore := srv.GetBlogListByTag(param.Name, param.PageNum, param.PageSize)
	c.RespSuccess(iris.Map{
		"list":    list,
		"hasMore": hasMore,
	}, "获取成功")
}

func (c *BlogController) GetAllTags() {
	srv := services.NewBlogService()
	list := srv.GetAllTags()
	c.RespSuccess(list, "获取成功")
}

func (c *BlogController) GetCateBlogList() {
	type Param struct {
		CateId   int64 `json:"cateId" validate:"required"`   //分类id
		PageNum  int   `json:"pageNum" validate:"required"`  //第几页，从1开始
		PageSize int   `json:"pageSize" validate:"required"` //每页的数量
	}

	var param Param
	if err := c.Ctx.ReadJSON(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	srv := services.NewBlogService()
	list, hasMore := srv.GetCateBlogList(param.CateId, param.PageNum, param.PageSize)
	c.RespSuccess(iris.Map{
		"list":    list,
		"hasMore": hasMore,
	}, "获取成功")

}

func (c *BlogController) GetBlogListByAction() {
	type Param struct {
		PageNum  int `json:"pageNum" validate:"required"`  //第几页，从1开始
		PageSize int `json:"pageSize" validate:"required"` //每页的数量
		Action   int `json:"action" validate:"oneof=0 1"`  //0：获取最新博客，1：获取最热博客
	}

	var param Param
	if err := c.Ctx.ReadJSON(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&param); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	srv := services.NewBlogService()
	list, hasMore := srv.GetBlogListByAction(param.Action, param.PageNum, param.PageSize)
	c.RespSuccess(iris.Map{
		"list":    list,
		"hasMore": hasMore,
	}, "获取成功")
}

func (c *BlogController) GetCateList() {
	srv := services.NewCategoryService()
	list := srv.GetCategoryList(models.CateStatusNormal)
	c.RespSuccess(list, "获取成功")
}

func (c *BlogController) GetBlogDetail() {
	var bean beans.GetBlogDetailBean
	if err := c.Ctx.ReadJSON(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	valid := validator.New()
	if err := valid.Struct(&bean); err != nil {
		c.RespFailedMessage("参数错误: " + err.Error())
		return
	}

	srv := services.NewBlogService()
	blog, list := srv.GetBlogDetail(bean.Id)
	if blog == nil || blog.Status != models.BlogStatusPublished {
		c.RespFailedMessage("博客不存在")
	} else {
		c.RespSuccess(iris.Map{
			"blog": blog,
			"list": list,
		}, "获取博客成功")
	}
}
