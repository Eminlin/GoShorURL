package model

import (
	"encoding/json"
	"net/http"
)

//ErrMsg
type ErrMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"errMsg"`
}

//SuccessAddRtn
type SuccessAddRtn struct {
	OriginURL string `json:"originURL"`
	ShortURL  string `json:"shortURL"`
	ShortKey  string `json:"shortKey"`
}

func CommonErrResp(w http.ResponseWriter, msg string) {
	w.WriteHeader(406)
	w.Header().Set("Content-Type", "application/json")
	m, _ := json.Marshal(ErrMsg{Code: 406, Msg: msg})
	w.Write(m)
}

func CommonSuccessResp(w http.ResponseWriter, v interface{}) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	m, _ := json.Marshal(v)
	w.Write(m)
}
