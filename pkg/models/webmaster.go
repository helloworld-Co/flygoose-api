package models

import "time"

// 站长信息表
type Webmaster struct {
	ID         int64           `json:"id"`         //id
	Intro      string          `json:"intro"`      //站长简介
	Slogan     string          `json:"slogan"`     //站长个性签名
	Nicker     string          `json:"nicker"`     //站长昵称
	Avatar     string          `json:"avatar"`     //站长头像
	Job        string          `json:"job"`        //站长职业
	Email      string          `json:"email"`      //站长邮箱
	QQ         string          `json:"qq"`         //站长QQ号
	Wechat     string          `json:"wechat"`     //站长微信二维码
	RewardCode string          `json:"rewardCode"` //站长打赏二维码
	Status     WebmasterStatus `json:"status"`     //状态
	UpdateTime time.Time       `json:"updateTime"` //更新时间
}

// 网站信息状态
type WebmasterStatus int

const (
	WebmasterStatusOffline = 0 //未启用
	WebmasterStatusOnline  = 1 //已启用
)
