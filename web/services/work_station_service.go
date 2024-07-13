package services

import (
	"flygoose/pkg/beans"
	"flygoose/web/daos"
)

type WorkStationService struct {
}

func NewWorkStationService() *WorkStationService {

	return &WorkStationService{}
}
func (*WorkStationService) GetStatistics() *beans.Statistics {
	var statistics beans.Statistics
	statistics.BlogCnt = daos.NewBlogDao().GetTotal()
	statistics.NewBlogTodayCnt = daos.NewBlogDao().GetTodayTotal()
	statistics.SpecialCnt = daos.NewSpecialDao().GetTotal()
	statistics.NewSpecialTodayCnt = daos.NewSpecialDao().GetTodayTotal()
	statistics.LinkCnt = daos.NewSpecialDao().GetTotal()
	statistics.AdminCnt = daos.NewAccessDao().GetTotal()
	statistics.BannerCnt = daos.NewBannerDao().GetTotal()
	return &statistics
}
