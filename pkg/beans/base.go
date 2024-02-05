package beans

type PageBean struct {
	PageNum  int `json:"pageNum" validate:"gt=0"`         // >0
	PageSize int `json:"pageSize" validate:"lte=50,gt=0"` // >0 && <=50
}

type PageStatusBean struct {
	PageBean
	Status int `json:"status" validate:"oneof=0 1 -1"`
}
