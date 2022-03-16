package server

import (
	"GoShortURL/common"
	"GoShortURL/model"
	"errors"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spaolacci/murmur3"
)

type Server struct {
	Log *common.Log
}

func NewServer(log *common.Log) *Server {
	return &Server{
		Log: log,
	}
}

func (s *Server) WebRun() {
	http.HandleFunc("/", index)
	http.HandleFunc("/manage", manage)
	http.HandleFunc("/manage/add", add)
	http.HandleFunc("/manage/del", del)

	s.Log.Infof("web start at http://127.0.0.1:%s", common.AppConf.App.APIPort)

	if err := http.ListenAndServe(":"+common.AppConf.App.APIPort, nil); err != nil {
		s.Log.Errorf("listen and serve err:%s\n", err.Error())
	}
}

//getMurmur
func getMurmur(text string) string {
	switch common.AppConf.App.MurmurBit {
	case 32:
		return common.Uint32ToB62(murmur3.Sum32([]byte(text)))
	case 64:
		return common.Uint64ToB62(murmur3.Sum64([]byte(text)))
	default:
		panic("Config MurmurBit can only be 32 or 64")
	}
}

//getMurmurWithSeed
func getMurmurWithSeed(text string) string {
	switch common.AppConf.App.MurmurBit {
	case 32:
		return common.Uint32ToB62(murmur3.Sum32WithSeed([]byte(text), uint32(time.Now().Unix())))
	case 64:
		return common.Uint64ToB62(murmur3.Sum64WithSeed([]byte(text), uint32(time.Now().Unix())))
	default:
		panic("Config MurmurBit can only be 32 or 64")
	}
}

//buildShortKey
func buildShortKey(maxTryTimes int, url string) (status bool, shortKey string, err error) {
	shortKey = getMurmur(url)
	for i := 0; i <= maxTryTimes; i++ {
		if common.RedisClient.HExists(common.AppConf.Redis.Key, shortKey).Val() {
			cmdRes, err := common.RedisClient.HMGet(common.AppConf.Redis.Key, shortKey).Result()
			if err != nil {
				return false, "", err
			}
			if cmdRes[0].(string) == url {
				return true, shortKey, nil
			}
			shortKey = getMurmurWithSeed(url)
			continue
		}
		if checkNotExistAndSet(shortKey) {
			return true, shortKey, nil
		}
		break
	}
	if shortKey == "" {
		return false, "", errors.New("buildShortKey maximum number of times")
	}
	return false, shortKey, nil
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
