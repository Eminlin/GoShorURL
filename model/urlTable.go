package model

//URL Record url
type URLTable struct {
	ShortKey    string `gorm:"short_key"`
	OriginURL   string `gorm:"origin_url"`
	Remark      string `gorm:"remark"`
	CreatedTime int64  `gorm:"created_time"`
	UpdateTime  int64  `gorm:"update_time"`
}
