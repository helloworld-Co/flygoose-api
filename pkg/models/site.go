package models

import "time"

// 网站信息表
type Site struct {
	ID         int64      `json:"id"`         //id
	Title      string     `json:"title"`      //网站名称
	Intro      string     `json:"intro"`      //网站简介
	Slogan     string     `json:"slogan"`     //个性签名
	Cover      string     `json:"cover"`      //网站背景图
	Copyright  string     `json:"copyright"`  //版权信息
	Icp        string     `json:"icp"`        //备案号
	UpdateTime time.Time  `json:"updateTime"` //更新时间
	Extra      string     `json:"extra"`      //额外的信息, json字符串,具体有哪些信息，暂未定
	Status     SiteStatus `json:"status"`     //网站信息状态
}

// 网站信息状态
type SiteStatus int

const (
	SiteStatusOffline = 0 //未启用
	SiteStatusOnline  = 1 //已启用
)
