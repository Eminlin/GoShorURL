package server

import (
	"GoShortURL/common"
	"GoShortURL/config"
	"fmt"
	"log"

	"github.com/spaolacci/murmur3"
)

func Run() {
	go WebRun()
	url := ""
	if !common.CheckURL(url) {
		log.Println("It seems doesn't look like a link")
		return
	}
	m := getMurmur(url)
	fmt.Println(m)

	select {}
}

//DupliCheck Duplicate URL data detection
func DupliCheck(url string) bool {
	return false
}

//getMurmur
func getMurmur(text string) string {
	switch config.MurmurBit {
	case 32:
		return common.Uint32ToB62(murmur3.Sum32([]byte(text)))
	case 64:
		return common.Uint64ToB62(murmur3.Sum64([]byte(text)))
	default:
		panic("Config MurmurBit can only be 32 or 64")
	}
}
