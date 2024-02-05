package beans

type Statistics struct {
	BlogCnt            int64 `json:"blogCnt"`            //博客数
	SpecialCnt         int64 `json:"specialCnt"`         //专栏数
	LinkCnt            int64 `json:"linkCnt"`            //友链数
	AdminCnt           int64 `json:"adminCnt"`           //管理员数
	BannerCnt          int64 `json:"bannerCnt"`          //轮播图数
	NewBlogTodayCnt    int64 `json:"newBlogToday"`       //今日新增博客数
	NewSpecialTodayCnt int64 `json:"newSpecialTodayCnt"` //今日新增专栏数
}
