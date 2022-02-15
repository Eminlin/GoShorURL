package common

import (
	"GoShortURL/model"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

//Config config
type Config struct {
}

type ConfigFile struct {
	FileName string
	Value    interface{}
}

func init() {
	c := NewConfig()
	if err := c.GetConfig(); err != nil {
		log.Fatalf("loading config err %s", err.Error())
	}
}

func NewConfig() *Config {
	return &Config{}
}

var pathSeparator = string(os.PathSeparator)

//AppConf app.conf struct
var AppConf model.App

//GetConfig GetConfig
func (c *Config) GetConfig() error {
	conf := []ConfigFile{}
	conf = append(conf, ConfigFile{"app.toml", &AppConf})
	configPath, err := getCurPath("GoShortURL")
	if err != nil {
		return err
	}
	return parse(conf, configPath)
}

//parse decode toml
func parse(conf []ConfigFile, path string) error {
	for _, conf := range conf {
		_, err := toml.DecodeFile(path+pathSeparator+conf.FileName, conf.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

//getCurPath get project path
func getCurPath(projectName string) (string, error) {
	curPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	curPath += pathSeparator
	index := strings.Index(curPath, pathSeparator+projectName+pathSeparator)
	if index >= 0 {
		curPath = subStr(curPath, 0, index) + pathSeparator + projectName + pathSeparator + "config" + pathSeparator
	}
	return curPath, nil
}

func subStr(str string, start, length int) string {
	runeStr := []rune(str)
	lenStr := len(runeStr)
	if start < 0 || length <= 0 {
		return ""
	}
	if start >= lenStr {
		return ""
	}
	end := start + length
	if end > lenStr-1 {
		end = lenStr
	}
	return string(runeStr[start:end])
}
