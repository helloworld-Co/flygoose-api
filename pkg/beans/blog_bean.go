package beans

type BlogPageStatusBean struct {
	PageBean
	Status int `json:"status" validate:"oneof=0 10 20 30"` //博客的状态值
}

type UpdateBlogStatusBean struct {
	Id     int64 `json:"id" validate:"gt=0"`               //博客id，必须大于0
	Status int   `json:"status" validate:"oneof=10 20 30"` //博客的状态值
}

type BlogPageWord struct {
	PageBean
	Word string `json:"word" validate:"required"` //关键字
}

type SearchBlogBean struct {
	PageBean
	Word   string `json:"word"`   //关键字
	Status int    `json:"status"` //博客的状态值
}

type GetBlogDetailBean struct {
	Id int64 `json:"id" validate:"gt=0"` //博客id，必须大于0
}

type CreateBlogBean struct {
	Title     string `json:"title" validate:"required"`
	Intro     string `json:"intro"`
	Content   string `json:"content"`
	Html      string `json:"html"`
	Tags      string `json:"tags"`
	Thumbnail string `json:"thumbnail"`
	IsHtml    int    `json:"isHtml" validate:"oneof=0 1"` //0：md编辑器产生的内容，1：富文本编辑器产生的内容
	IsTop     int    `json:"isTop" validate:"oneof=0 1"`  //是否置顶, 0:不置顶， 1：置顶
	CateId    int64  `json:"cateId" validate:"gt=0"`      //博客分类id
}

type UpdateBlogBean struct {
	Title   string `json:"title"`
	Intro   string `json:"intro"`
	Content string `json:"content"`
	Html    string `json:"html"`
	Tags    string `json:"tags"`
	IsHtml  int    `json:"isHtml"` //0：md编辑器产生的内容，1：富文本编辑器产生的内容
	IsTop   int    `json:"isTop"`  //是否置顶, 0:不置顶， 1：置顶
	CateId  int64  `json:"cateId"` //博客分类id
}
