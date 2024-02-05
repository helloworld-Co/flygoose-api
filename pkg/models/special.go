package models

import "time"

// Special 专栏表
type Special struct {
	ID          int64         `json:"id"`          //id
	Title       string        `json:"title"`       //标题
	Intro       string        `json:"intro"`       //简介
	Cover       string        `json:"cover"`       //封面
	CreateTime  time.Time     `json:"createTime"`  //创建时间
	UpdateTime  time.Time     `json:"updateTime"`  //更新时间
	PublishTime time.Time     `json:"publishTime"` //发布时间
	Status      SpecialStatus `json:"status"`      //状态
}

// SpecialStatus 专栏状态
type SpecialStatus int

const (
	SpecialStatusCreated   SpecialStatus = 10 //未发布(新建状态,或者已下架)
	SpecialStatusDeleted   SpecialStatus = 20 //已删除状态
	SpecialStatusPublished SpecialStatus = 30 //已发布状态
)
