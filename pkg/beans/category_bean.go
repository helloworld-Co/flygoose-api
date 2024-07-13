package beans

type CatePageStatusBean struct {
	Status int `json:"status" validate:"oneof=0 1 -1"` // -1为全部
}
