package server

import (
	"GoShortURL/common"
	"GoShortURL/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"unicode"

	"github.com/jinzhu/gorm"
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
	path := r.URL.Path[1:]
	for _, v := range path {
		if !unicode.IsLetter(v) && !unicode.IsDigit(v) || unicode.Is(unicode.Han, v) {
			model.CommonErrResp(w, "invalid path")
			return
		}
	}
	go log.Infof("Path:%s Method:%s UserAgent:%v RemoteAddr:%s", r.URL.Path, r.Method, r.UserAgent(), r.RemoteAddr)
	if !common.RedisClient.HExists(common.AppConf.Redis.Key, path).Val() {
		if s := checkNotExistAndSet(path); !s {
			notFound(w)
			return
		}
	}
	cmdRes, err := common.RedisClient.HMGet(common.AppConf.Redis.Key, path).Result()
	if err != nil {
		log.Errorln(err)
		model.CommonErrResp(w, fmt.Sprintf("system err:%s", err.Error()))
		return
	}
	w.Header().Set("Location", cmdRes[0].(string))
	w.Header().Set("Referer", common.AppConf.App.Host+r.URL.Path)
	w.WriteHeader(http.StatusFound) //302
	err = common.RedisClient.HIncrBy(common.AppConf.Redis.Key, path+"_visit", 1).Err()
	if err != nil {
		log.Errorln(err)
		return
	}
}

//checkNotExistAndSet check short key when not found in redis
func checkNotExistAndSet(path string) bool {
	var urlTable model.URLTable
	err := common.DB.Table(common.AppConf.MySQL.URLTable).
		Select("origin_url").Where("short_key = ?", path).
		Find(&urlTable).Limit(1).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return false
		}
		log.Errorln(err.Error())
		return false
	}
	if urlTable.OriginURL == "" {
		return false
	}
	boolCmd := common.RedisClient.HSet(common.AppConf.Redis.Key, path, urlTable.OriginURL)
	if boolCmd.Err() != nil {
		log.Errorln(boolCmd.Err())
		return false
	}
	boolCmd = common.RedisClient.HSet(common.AppConf.Redis.Key, path+"_visit", 0)
	if boolCmd.Err() != nil {
		log.Errorln(boolCmd.Err())
		return false
	}
	return true
}

//add new short url
func add(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln(err)
		model.CommonErrResp(w, fmt.Sprintf("system err:%s", err.Error()))
		return
	}
	if len(b) == 0 {
		model.CommonErrResp(w, "invalid param empty")
		return
	}
	param := model.AddParam{}
	if err = json.Unmarshal(b, &param); err != nil {
		log.Errorln(err)
		model.CommonErrResp(w, fmt.Sprintf("system err:%s", err.Error()))
		return
	}
	if !common.CheckURL(param.URL) {
		model.CommonErrResp(w, "invalid param url")
		return
	}
	status, short, err := buildShortKey(6, param.URL)
	if err != nil {
		log.Errorln(err.Error())
		model.CommonErrResp(w, fmt.Sprintf("system err:%s", err.Error()))
		return
	}
	if status {
		model.CommonSuccessResp(w, model.SuccessAddRtn{
			OriginURL: param.URL,
			ShortURL:  common.AppConf.App.Host + "/" + short,
			ShortKey:  short,
		})
		return
	}
	pipe := common.RedisClient.TxPipeline()
	defer pipe.Close()
	boolCmd := pipe.HSet(common.AppConf.Redis.Key, short, param.URL)
	if boolCmd.Err() != nil {
		log.Errorln(boolCmd.Err())
		model.CommonErrResp(w, fmt.Sprintf("system err:%s", boolCmd.Err()))
		return
	}
	boolCmd = pipe.HSet(common.AppConf.Redis.Key, short+"_visit", 0)
	if boolCmd.Err() != nil {
		log.Errorln(boolCmd.Err())
		model.CommonErrResp(w, fmt.Sprintf("system err:%s", boolCmd.Err()))
		return
	}
	if err = common.DB.Table(common.AppConf.MySQL.URLTable).Create(model.URLTable{
		ShortKey:    short,
		OriginURL:   param.URL,
		Remark:      param.Remark,
		CreatedTime: time.Now().Unix(),
		UpdateTime:  0,
	}).Error; err != nil {
		log.Errorln(err)
		pipe.Discard()
		model.CommonErrResp(w, fmt.Sprintf("system err:%s", err.Error()))
		return
	}
	if _, err = pipe.Exec(); err != nil {
		log.Errorln(err)
		pipe.Discard()
		model.CommonErrResp(w, fmt.Sprintf("system err:%s", err.Error()))
		return
	}
	model.CommonSuccessResp(w, model.SuccessAddRtn{
		OriginURL: param.URL,
		ShortURL:  common.AppConf.App.Host + "/" + short,
		ShortKey:  short,
	})
	log.Infoln(short, param.URL)
}

//add new short url
func del(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln(err)
		model.CommonErrResp(w, fmt.Sprintf("system err:%s", err.Error()))
		return
	}
	if len(b) == 0 {
		model.CommonErrResp(w, "invalid param empty")
		return
	}
	delParam := model.DelParam{}
	if err = json.Unmarshal(b, &delParam); err != nil {
		log.Errorln(err)
		model.CommonErrResp(w, fmt.Sprintf("system err:%s", err.Error()))
		return
	}
	if delParam.ShortKey == "" {
		model.CommonErrResp(w, "invalid param shortKey empty")
		return
	}
	if !common.RedisClient.HExists(common.AppConf.Redis.Key, delParam.ShortKey).Val() {
		model.CommonErrResp(w, "not found")
		return
	}
	pipe := common.RedisClient.TxPipeline()
	defer pipe.Close()
	if err := pipe.HDel(common.AppConf.Redis.Key, delParam.ShortKey).Err(); err != nil {
		log.Errorln(err)
		pipe.Discard()
		model.CommonErrResp(w, "del err")
		return
	}
	if err := pipe.HDel(common.AppConf.Redis.Key, delParam.ShortKey+"_visit").Err(); err != nil {
		log.Errorln(err)
		pipe.Discard()
		model.CommonErrResp(w, "del err")
		return
	}
	if _, err = pipe.Exec(); err != nil {
		log.Errorln(err)
		pipe.Discard()
		model.CommonErrResp(w, "del err")
		return
	}
	model.CommonSuccessResp(w, "success")
}

//notFound 404 page
func notFound(w http.ResponseWriter) {
	if !common.CheckURL(common.AppConf.App.NotFoundPage) {
		model.CommonErrResp(w, "not found")
		return
	}
	w.Header().Set("Location", common.AppConf.App.NotFoundPage)
	w.WriteHeader(http.StatusFound) //302
}

//manage manage page
func manage(w http.ResponseWriter, r *http.Request) {

}
