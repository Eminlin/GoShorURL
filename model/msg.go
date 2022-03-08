package model

import (
	"encoding/json"
	"net/http"
)

type ErrMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"errMsg"`
}

func CommonErrResp(w http.ResponseWriter, msg string) {
	w.WriteHeader(406)
	m, _ := json.Marshal(ErrMsg{Code: 406, Msg: msg})
	w.Write(m)
}
