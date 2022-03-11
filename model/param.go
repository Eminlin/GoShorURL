package model

type AddParam struct {
	URL string `json:"url"`
}

type RedisHSet struct {
	OriginURL string `json:"originURL"`
	ShortKey  string `json:"shortKey"`
	Visit     string `json:"visit"`
}

type DelParam struct {
	ShortKey string `json:"shortKey"`
}
