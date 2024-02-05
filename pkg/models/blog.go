package models

import "time"

// 博客表
type Blog struct {
	ID          int64      `json:"id"`          //博客id
	Title       string     `json:"title"`       //博客的标题
	Intro       string     `json:"intro"`       //博客的简介
	Content     string     `json:"content"`     //md编辑器产生的内容时，对应的md的内容
	Html        string     `json:"html"`        //md编辑器产生的内容时，对应的生成的html
	IsHtml      int        `json:"isHtml"`      //0：md编辑器产生的内容，1：富文本编辑器产生的内容
	IsTop       int        `json:"isTop"`       //是否置顶, 0:不置顶， 1：置顶
	ReadCount   int        `json:"readCount"`   //阅读数
	CreateTime  time.Time  `json:"createTime"`  //创建时间
	UpdateTime  time.Time  `json:"updateTime"`  //更新时间
	PublishTime time.Time  `json:"publishTime"` //发布时间
	Thumbnail   string     `json:"thumbnail"`   //缩略图,是一个url
	Status      BlogStatus `json:"status"`      //博客状态
	CateId      int64      `json:"cateId"`      //博客分类id
	Tags        string     `json:"tags"`        //标签，英文逗号,分隔,如："并发,线程,Java"
}

// 博客状态
type BlogStatus int

const (
	BlogStatusCreated   BlogStatus = 10 //未发布(新建状态,或者已下架)
	BlogStatusDeleted   BlogStatus = 20 //已删除
	BlogStatusPublished BlogStatus = 30 //已发布
)
