package models

import "time"

// 轮播图表
type Banner struct {
	ID         int64        `json:"id"`         //banner id
	Title      string       `json:"title"`      //banner 标题
	Url        string       `json:"url"`        //banner url
	TargetUrl  string       `json:"targetUrl"`  //banner 跳转到的目标url
	Seq        int          `json:"seq"`        //序号，从小到大排序
	CreateTime time.Time    `json:"createTime"` //创建时间
	Status     BannerStatus `json:"status"`     //banner 状态
}

// banner状态
type BannerStatus int

const (
	BannerStatusOffline LinkStatus = 0 //已下架
	BannerStatusNormal  LinkStatus = 1 //正常
)
