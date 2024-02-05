package services

import (
	"errors"
	"flygoose/pkg/models"
	"flygoose/pkg/tools"
	"flygoose/web/daos"
	"time"
)

type AccessService struct {
	accessDao *daos.AccessDao
}

func NewAccessService() *AccessService {
	return &AccessService{accessDao: daos.NewAccessDao()}
}

func (s *AccessService) LoginIn(username string, password string) (bool, string) {
	//查找
	admin, err := s.accessDao.FirstByUsernameAndPassword(username, password)
	if err != nil {
		return false, err.Error()
	}

	if admin.Status == models.AdminStatusBlocked {
		return false, "此账号已被封禁,请联系站长解封"
	}

	//更新
	fields := make([]string, 0)
	var m models.Admin

	m.LoginTime = time.Now()
	m.ValidTime = time.Now().Add(time.Hour * 24 * 7) //7天有效期
	m.Token = tools.GenToken()

	fields = append(fields, "LoginTime")
	fields = append(fields, "ValidTime")
	fields = append(fields, "Token")

	err = s.accessDao.Update(admin.ID, fields, &m)
	if err != nil {
		return false, "更新数据失败"
	}

	return true, m.Token
}

func (s *AccessService) Logout(uid int64) error {
	admin := s.accessDao.FindAdminByUid(uid)
	if admin == nil {
		return errors.New("没有此用户")
	}

	//更新
	fields := make([]string, 0)
	var m models.Admin

	m.ValidTime = time.Now()
	m.Token = ""

	fields = append(fields, "ValidTime")
	fields = append(fields, "Token")

	return s.accessDao.Update(admin.ID, fields, &m)
}

func (s *AccessService) FirstAdminByToken(token string) (*models.Admin, error) {
	return s.accessDao.FirstAdminByToken(token)
}
