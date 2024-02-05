package beans

type CreateBannerBean struct {
	Title     string `json:"title" validate:"required"`
	Url       string `json:"url" validate:"required"`
	TargetUrl string `json:"targetUrl"`
	Seq       int    `json:"seq" validate:"required"`
}

type BannerPageStatusBean struct {
	PageBean
	Status int `json:"status" validate:"oneof=0 1 -1"` //0:已下架， 1：已上架， -1：全部
}
