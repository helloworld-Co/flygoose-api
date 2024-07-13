package daos

import (
	"flygoose/datasource"
	"flygoose/pkg/models"
	"flygoose/pkg/tlog"
	"time"

	"gorm.io/gorm"
)

type AccessDao struct {
	db *gorm.DB
}

func NewAccessDao() *AccessDao {
	return &AccessDao{db: datasource.GetMasterDB()}
}

func (dao *AccessDao) FirstByUsernameAndPassword(username string, password string) (*models.Admin, error) {
	var admin models.Admin
	result := dao.db.Where("phone=? and password=?", username, password).First(&admin)
	if result.Error != nil {
		tlog.Error2("AccessDao:FirstByUsernameAndPassword出错", result.Error)
		return nil, result.Error
	}
	return &admin, nil
}

func (dao *AccessDao) CountUsername(username string) (int64, error) {
	var cnt int64
	result := dao.db.Model(&models.Admin{}).Where("phone=? and status = 1", username).Count(&cnt)
	if result.Error != nil {
		tlog.Error2("查询初始化默认用户名失败", result.Error)
		return -1, result.Error
	}
	return cnt, nil
}

func (dao *AccessDao) Create(m *models.Admin) error {
	result := dao.db.Create(m)
	if result.Error != nil {
		tlog.Error2("AccessDao:Create出错", result.Error)
	}
	return result.Error
}

func (dao *AccessDao) Update(id int64, fields []string, admin *models.Admin) error {
	result := dao.db.Model(&models.Admin{}).Select(fields).Where("id=?", id).Updates(admin)
	if result.Error != nil {
		tlog.Error2("AccessDao:Update出错", result.Error)
	}
	return result.Error
}

func (dao *AccessDao) FindAminWithToken(token string) *models.Admin {
	var admin models.Admin
	result := dao.db.First(&admin, "token=?", token)
	if result.Error != nil {
		tlog.Error2("AccessDao:FindAminWithToken", result.Error)
		return nil
	}

	return &admin
}

func (dao *AccessDao) FindAdminByUid(uid int64) *models.Admin {
	var admin models.Admin
	result := dao.db.First(&admin, "id=?", uid)
	if result.Error != nil {
		tlog.Error2("AccessDao:FindAdminByUid", result.Error)
		return nil
	}

	return &admin
}

func (dao *AccessDao) CleanUserLoginInfo(uid int64) {
	result := dao.db.Model(&models.Admin{}).Where("id=?", uid).Updates(map[string]interface{}{"token": "", "valid_time": time.Now()})
	if result.Error != nil {
		tlog.Error2("AccessDao:CleanUserLoginInfo", result.Error)
	}
}

func (dao *AccessDao) FirstAdminByToken(token string) (*models.Admin, error) {
	var admin models.Admin
	result := dao.db.First(&admin, "token=?", token)
	if result.Error != nil {
		tlog.Error2("AccessDao:FindAdminByUid", result.Error)
		return nil, result.Error
	}

	return &admin, nil
}

func (dao *AccessDao) GetTotal() int64 {
	var total int64
	result := dao.db.Model(&models.Admin{}).Where("status = ?", models.AdminStatusNormal).Count(&total)
	if result.Error != nil {
		tlog.Error2("AccessDao GetTotal报错", result.Error)
		return 0
	}

	return total
}
