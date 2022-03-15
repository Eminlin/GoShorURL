package model

type AddParam struct {
	URL    string `json:"url"`
	Remark string `json:"remark"`
}

type RedisHSet struct {
	OriginURL string `json:"originURL"`
	ShortKey  string `json:"shortKey"`
	Visit     string `json:"visit"`
}

type DelParam struct {
	ShortKey string `json:"shortKey"`
}
