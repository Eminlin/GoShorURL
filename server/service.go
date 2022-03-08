package server

import (
	"GoShortURL/common"
	"net/http"

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

func (s *Server) Run() {
	// url := ""
	// if !common.CheckURL(url) {
	// 	s.Log.Errorf("It seems doesn't look like a link:%s\n", url)
	// 	return
	// }
	// m := getMurmur(url)
	// fmt.Println(m)
	s.WebRun()
}

func (s *Server) WebRun() {
	http.HandleFunc("/", index)
	http.HandleFunc("/manage", manage)
	http.HandleFunc("/manage/add", add)

	s.Log.Infof("web start at http://127.0.0.1:%s", common.AppConf.App.APIPort)

	if err := http.ListenAndServe(":"+common.AppConf.App.APIPort, nil); err != nil {
		s.Log.Errorf("listen and serve err:%s\n", err.Error())
	}
}

//DupliCheck Duplicate URL data detection
func DupliCheck(url string) bool {
	return false
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
