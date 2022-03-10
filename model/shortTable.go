package model

//URL Record url
type URLTable struct {
	ID          int64
	ShortKey    string
	OriginURL   string
	Remark      string
	CreatedTime int64
	UpdateTime  int64
}
