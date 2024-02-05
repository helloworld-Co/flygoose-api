package beans

type CreateLinkBean struct {
	Title     string `json:"title" validate:"required"`
	Url       string `json:"url" validate:"required"`
	Seq       int    `json:"seq"`
	ValidTime int64  `json:"validTime"`
	Remark    string `json:"remark"`
}

type LinkPageStatusBean struct {
	PageBean
	Status int `json:"status" validate:"oneof=0 1 -1"` //0:已下架， 1：已上架， -1：全部
}
