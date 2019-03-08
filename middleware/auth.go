package middleware

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/bingtianbaihua/goproxy/config"
	"github.com/bingtianbaihua/goproxy/log"
)

/*
	https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/407

	HTTP/1.1 407 Proxy Authentication Required
	Date: Wed, 21 Oct 2015 07:28:00 GMT
	Proxy-Authenticate: Basic realm="Access to internal site"
*/

var (
	proxyAuthError = []byte("HTTP/1.1 407 Proxy Authorization Required\r\nProxy-Authenticate: Basic realm=\"Access to internal site\"\r\n\r\n")
)

type AuthAdapter struct {
	IsAuth bool
	User   map[string]string
}

func NewAuthAdapter(cfg *config.Config) (*AuthAdapter, error) {
	return &AuthAdapter{
		IsAuth: cfg.IsAuth,
		User:   cfg.User,
	}, nil
}

func (p *AuthAdapter) HandleTask(w http.ResponseWriter, r *http.Request, next func(http.ResponseWriter, *http.Request)) {
	if p.IsAuth == true {
		if name, err := p.auth(w, r); err != nil {
			log.Warn("%s can not successfully access %v", name, err)
			return
		}
	}
	next(w, r)
}

func (p *AuthAdapter) auth(rw http.ResponseWriter, req *http.Request) (string, error) {
	auth := req.Header.Get("Proxy-Authorization")
	auth = strings.Replace(auth, "Basic ", "", 1)

	if auth == "" {
		writeResp(rw, proxyAuthError)
		return "", errors.New("Need Proxy Authorization!")
	}

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Proxy-Authorization
	data, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		log.Debug("when decoding %v, got an error of %v", auth, err)
		return "", errors.New("Fail to decoding Proxy-Authorization")
	}

	var user, passwd string

	// username:password
	UserPasswdPair := strings.Split(string(data), ":")
	if len(UserPasswdPair) != 2 {
		writeResp(rw, proxyAuthError)
		return "", errors.New("Fail to log in")
	} else {
		user = UserPasswdPair[0]
		passwd = UserPasswdPair[1]
	}

	if p.checkUserPwd(user, passwd) == false {
		writeResp(rw, proxyAuthError)
		return "", errors.New("Fail to log in")
	}
	return user, nil
}

func (p *AuthAdapter) checkUserPwd(user, passwd string) bool {
	passwd = strings.Trim(passwd, "")
	if user != "" && passwd != "" && p.User[user] == passwd {
		return true
	}
	return false
}

func writeResp(rw http.ResponseWriter, data []byte) error {
	n, err := rw.Write(data)
	if err != nil {
		log.Error("write data error: %v", err)
		return errors.New("InternalServerError")
	}
	if n != len(data) {
		return errors.New("short written")
	}
	return nil
}
