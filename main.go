package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"rev-proxy/util"
)

func main(){
	http.ListenAndServe(":8080", &webHandler{})
}



type webHandler struct {

}

func (wh *webHandler) ServeHTTP(w http.ResponseWriter,r *http.Request) {

	defer func() {
		if err:=recover();err!=nil{
			w.WriteHeader(500)
			log.Println(err)
		}
	}()

	lb:=&util.LoadBalance{}


	for k,v:=range util.Config.Mapping{
		if matched, _ := regexp.MatchString(k, r.URL.Path);matched==true{
			//go 内置的反向代理功能
			selected:=lb.SelectAlg(util.Config.LoadBalance,v,r.RemoteAddr)
			log.Println(selected)
			parse, _ := url.Parse(selected.Host)
			proxy:= httputil.NewSingleHostReverseProxy(parse)
			proxy.ServeHTTP(w,r)

			//util.ProxyUrl(w,r,v)
			return
		}
	}
	w.Write([]byte("<h1>404</404>"))
}