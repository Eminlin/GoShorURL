package server

import (
	"GoShortURL/common"
	"GoShortURL/model"
	"net/http"
	"unicode"
)

var log common.Log

//index short url request redict
func index(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" || r.RequestURI == "/favicon.ico" {
		model.CommonErrResp(w, "path param is empty")
		return
	}
	if r.Method != "GET" {
		model.CommonErrResp(w, "invalid request method")
		return
	}
	path := r.URL.Path
	path = path[1:]
	for _, v := range path {
		if !unicode.IsLetter(v) && !unicode.IsDigit(v) || unicode.Is(unicode.Han, v) {
			model.CommonErrResp(w, "invalid path")
			return
		}
	}
	go log.Infof("Path:%s Method:%s UserAgent:%v RemoteAddr:%s", r.URL.Path, r.Method, r.UserAgent(), r.RemoteAddr)

	if !common.RedisClient.HExists(common.AppConf.Redis.Key, path).Val() {
		return
	}

	cmdRes, err := common.RedisClient.HMGet(common.AppConf.Redis.Key+path, "OriginURL", "Visit").Result()
	if err != nil {
		return
	}

	w.Header().Set("Location", cmdRes[0].(string))
	w.Header().Set("Referer", common.AppConf.App.Host+r.URL.Path)
	w.WriteHeader(302)

	cmdRes[1] = cmdRes[1].(int) + 1
	_, err = common.RedisClient.HMSet(common.AppConf.Redis.Key+path, map[string]interface{}{"Visit": cmdRes[1]}).Result()
	if err != nil {
		log.Errorln(err)
		return
	}
}

//manage manage page
func manage(w http.ResponseWriter, r *http.Request) {

}

//add new short url
func add(w http.ResponseWriter, r *http.Request) {
	if !common.CheckURL("") {
		return
	}
	m := map[string]interface{}{
		"OriginURL": "",
		"Visit":     0,
	}
	boolCmd := common.RedisClient.HMSet(
		common.AppConf.Redis.Key+(""),
		m,
	)
	if boolCmd.Err() != nil {
		log.Errorln(boolCmd.Err())
		return
	}

}

//add new short url
func del(w http.ResponseWriter, r *http.Request) {

}
