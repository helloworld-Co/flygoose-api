package beans

import "flygoose/pkg/models"

type CreateSpecialParam struct {
	Title string `json:"title" validate:"required"`
	Intro string `json:"intro"`
	Cover string `json:"cover"`
}

type AddSectionParam struct {
	SpecialId int64  `json:"specialId" validate:"required"`
	Title     string `json:"title" validate:"required"`
	Seq       int    `json:"seq" validate:"required"`
}

type SearchSpecialParam struct {
	PageBean
	Word   string `json:"word"`   //关键字
	Status int    `json:"status"` //博客的状态值
}

type SpecialBean struct {
	models.Special
	TotalCount     int `json:"totalCount"`     //总的小节数
	PublishedCount int `json:"publishedCount"` //已发布的小节数
}

type SpecialHomeBean struct {
	models.Special
	PublishedCount int64 `json:"publishedCount"` //已发布的小节数
	ReadCount      int64 `json:"readCount"`      //章节的阅读数
}

type GetSectionParam struct {
	SectionId int64 `json:"sectionId" validate:"required"`
}
