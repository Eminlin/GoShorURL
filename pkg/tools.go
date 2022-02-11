package pkg

import "net/url"

//CheckURL Check link is legal
func CheckURL(l string) bool {
	if _, err := url.ParseRequestURI(l); err != nil {
		return false
	}
	return true
}
