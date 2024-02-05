package models

import "time"

// 管理员表
type Admin struct {
	ID         int64       `json:"id"`              //id
	Phone      string      `json:"phone"`           //手机号
	Password   string      `json:"password"`        //密码
	Nicker     string      `json:"nicker"`          //昵称
	Avatar     string      `json:"avatar"`          //头像
	Token      string      `json:"token,omitempty"` //登录后的token
	CreateTime time.Time   `json:"createTime"`      //账号创建时间
	ValidTime  time.Time   `json:"validTime"`       //登录后有效期
	LoginTime  time.Time   `json:"loginTime"`       //登录时间
	Status     AdminStatus `json:"status"`          //账号状态
}

// AdminStatus 管理员状态,默认为正常
type AdminStatus int

const (
	AdminStatusBlocked AdminStatus = 0 //封禁状态
	AdminStatusNormal  AdminStatus = 1 //正常可用的状态
)
