package models

import "time"

// 专栏中的小节表
type Section struct {
	ID          int64         `json:"id"`          //小节id
	Title       string        `json:"title"`       //小节的标题
	Intro       string        `json:"intro"`       //小节的简介，暂时用不到
	Content     string        `json:"content"`     //md编辑器产生的内容时，对应的md的内容
	Html        string        `json:"html"`        //md编辑器产生的内容时，对应的生成的html
	Tags        string        `json:"tags"`        //标签，英文逗号,分隔,如："并发,线程,Java"
	IsHtml      int           `json:"isHtml"`      //0：md编辑器产生的内容，1：富文本编辑器产生的内容
	ReadCount   int64         `json:"readCount"`   //阅读数
	SpecialId   int64         `json:"specialId"`   //所属专栏id
	Seq         int           `json:"seq"`         //序号，从小到大排序，默认100
	CreateTime  time.Time     `json:"createTime"`  //创建时间
	UpdateTime  time.Time     `json:"updateTime"`  //更新时间
	PublishTime time.Time     `json:"publishTime"` //发布时间
	Status      SectionStatus `json:"status"`      //小节状态
}

// 博客状态
type SectionStatus int

const (
	SectionStatusCreated   SectionStatus = 10 //未发布(新建状态,或者已下架)
	SectionStatusDeleted   SectionStatus = 20 //已删除
	SectionStatusPublished SectionStatus = 30 //已发布
)
