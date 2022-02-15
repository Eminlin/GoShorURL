package server

import (
	"GoShortURL/common"
	"net/http"
)

func WebRun() {
	http.Handle("/", http.NotFoundHandler())
	http.ListenAndServe(":"+common.AppConf.App.APIPort, nil)
}
