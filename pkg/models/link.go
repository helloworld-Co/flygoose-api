package models

import "time"

// 友链表
type Link struct {
	ID         int64      `json:"id"`         //id
	Title      string     `json:"title"`      //标题
	Url        string     `json:"url"`        //链接到的url
	Seq        int        `json:"seq"`        //序号，从小到大排序
	CreateTime time.Time  `json:"createTime"` //创建时间
	UpdateTime time.Time  `json:"updateTime"` //修改时间
	ValidTime  time.Time  `json:"validTime"`  //有效时间
	Remark     string     `json:"remark"`     //备注
	Status     LinkStatus `json:"status"`     //状态
}

// 友链状态
type LinkStatus int

const (
	LinkStatusOffline LinkStatus = 0 //已下架
	LinkStatusNormal  LinkStatus = 1 //正常
)
