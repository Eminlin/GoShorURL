package server

import (
	"GoShortURL/common"
	"net/http"
	"unicode"
)

var log common.Log

//index short url request redict
func index(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" || r.RequestURI == "/favicon.ico" {
		return
	}
	path := r.URL.Path
	path = path[1:]
	for _, v := range path {
		if !unicode.IsLetter(v) && !unicode.IsDigit(v) || unicode.Is(unicode.Han, v) {
			return
		}
	}
	go log.Infof("Path:%s Method:%s UserAgent:%v RemoteAddr:%s", r.URL.Path, r.Method, r.UserAgent(), r.RemoteAddr)

	w.Header().Set("Location", "")
	w.Header().Set("Referer", common.AppConf.App.Host+r.URL.Path)
	w.WriteHeader(302)
}

//manage manage page
func manage(w http.ResponseWriter, r *http.Request) {

}

//add new short url
func add(w http.ResponseWriter, r *http.Request) {

}

//add new short url
func del(w http.ResponseWriter, r *http.Request) {

}
