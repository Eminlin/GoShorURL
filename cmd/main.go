package main

import (
	"GoShortURL/common"
	"GoShortURL/server"
)

func init() {

}

func main() {
	server.NewServer(common.NewLog()).Run()
}
