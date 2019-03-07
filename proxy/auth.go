package proxy

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/bingtianbaihua/goproxy/log"
)

/*
	https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/407

	HTTP/1.1 407 Proxy Authentication Required
	Date: Wed, 21 Oct 2015 07:28:00 GMT
	Proxy-Authenticate: Basic realm="Access to internal site"
*/

var HTTP_407 = []byte("HTTP/1.1 407 Proxy Authorization Required\r\nProxy-Authenticate: Basic realm=\"Access to internal site\"\r\n\r\n")

func (proxy *ProxyServer) Auth(rw http.ResponseWriter, req *http.Request) bool {
	var err error
	if cnfg.Auth == true {
		if proxy.Name, err = proxy.auth(rw, req); err != nil {
			log.Debug("%s can not successfully access %v", proxy.Name, err)
			return true
		}
	} else {
		proxy.Name = "default-proxy"
	}

	log.Info("%s successfully log in!", proxy.Name)
	return false
}

func (proxy *ProxyServer) auth(rw http.ResponseWriter, req *http.Request) (string, error) {
	auth := req.Header.Get("Proxy-Authorization")
	auth = strings.Replace(auth, "Basic ", "", 1)

	if auth == "" {
		writeResp(rw, HTTP_407)
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
		writeResp(rw, HTTP_407)
		return "", errors.New("Fail to log in")
	} else {
		user = UserPasswdPair[0]
		passwd = UserPasswdPair[1]
	}

	if check(user, passwd) == false {
		writeResp(rw, HTTP_407)
		return "", errors.New("Fail to log in")
	}
	return user, nil
}

func writeResp(rw http.ResponseWriter, data []byte) error {
	_, err := rw.Write(data)
	if err != nil {
		log.Error("fail to write data to response")
		return errors.New("InternalServerError")
	}
	return nil
}

func check(User, passwd string) bool {
	if User != "" && passwd != "" && cnfg.User[User] == passwd {
		return true
	}
	return false
}
