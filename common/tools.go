package common

import "net/url"

//CheckURL Check link is legal
func CheckURL(l string) bool {
	if l == "" {
		return false
	}
	if _, err := url.ParseRequestURI(l); err != nil {
		return false
	}
	return true
}
