package helper

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pajri/personal-backend/config"
)

type CookieHelper struct{}

func (ch CookieHelper) SetHttpOnlyCookie(name, value string, expire time.Time) *http.Cookie {
	var host string

	//extract host
	u, err := url.Parse(config.Config.Host)
	if err != nil {
		log.Println("[SHO00] url parse error : " + err.Error())
	}
	hostSplit, _, _ := net.SplitHostPort(u.Host)
	if hostSplit != "" {
		host = hostSplit
	} else {
		host = u.Host
	}

	//set cookie
	cookie := &http.Cookie{}
	cookie.Name = name
	cookie.Value = value
	cookie.HttpOnly = true
	cookie.Domain = host
	cookie.Path = "/"
	cookie.Expires = expire

	return cookie
}

func (ch CookieHelper) RemoveHttpOnlyCookie(name string) *http.Cookie {
	var host string

	//extract host
	u, err := url.Parse(config.Config.Host)
	if err != nil {
		log.Println("[RHO00] url parse error : " + err.Error())
	}
	hostSplit, _, _ := net.SplitHostPort(u.Host)
	if hostSplit != "" {
		host = hostSplit
	} else {
		host = u.Host
	}

	//set cookie
	cookie := &http.Cookie{}
	cookie.Name = name
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Domain = host
	cookie.Path = "/"
	cookie.Expires = time.Time{}
	cookie.MaxAge = -1

	return cookie
}
