package models

import "time"

// 博客分类
type Category struct {
	ID         int64      `json:"id"`         //分类id
	Name       string     `json:"name"`       //名称
	Seq        int        `json:"seq"`        //序号，从小到大排序
	Icon       string     `json:"icon"`       //分类图标icon , url
	Font       string     `json:"font"`       //字体	，iconfont上的矢量字体
	Color      string     `json:"color"`      //分类颜色，iconfont上的矢量字体颜色，如：#aa33ff
	CreateTime time.Time  `json:"createTime"` //分类创建时间
	UpdateTime time.Time  `json:"updateTime"` //分类更新时间
	Status     CateStatus `json:"status"`     //状态
}

// 分类状态
type CateStatus int

const (
	CateStatusOffline = 0 //下架
	CateStatusNormal  = 1 //正常
)
