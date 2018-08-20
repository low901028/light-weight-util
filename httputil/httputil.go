package httputil

import (
	"net/http"
	"io"
	"io/ioutil"
	"net"
)

// GracefulClose获取数据直到EOF，并关闭
// 同时能够防止TCP/TLS关闭，connection能够被重用
func GracefulClose(resp *http.Response){
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}


// GetHostname返回reques host字段的hostname
// 若是host字段包含非法内容，hostname为空
func GetHostname(req *http.Request) string{
	if req == nil{
		return ""
	}
	h, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		return req.Host
	}
	return h
}
