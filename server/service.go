package server

import (
	"GoShortURL/common"
	"fmt"

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
	url := ""
	if !common.CheckURL(url) {
		s.Log.Fatalln("It seems doesn't look like a link")
		return
	}
	m := getMurmur(url)
	fmt.Println(m)
	WebRun()
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
