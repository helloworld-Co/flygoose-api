package models

import "time"

// 网站公告
type Notice struct {
	ID         int64     `json:"id"`         //id
	Title      string    `json:"title"`      //公告标题
	Content    string    `json:"content"`    //公告内容
	CreateTime time.Time `json:"createTime"` //创建时间
	UpdateTime time.Time `json:"updateTime"` //更新时间
	ValidTime  time.Time `json:"validTime"`  //有效时间
	Status     int       `json:"status"`     //公告状态1-上架0-下架
}
