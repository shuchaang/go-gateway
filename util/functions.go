package util

import (
	"io/ioutil"
	"log"
	"net/http"
)

func CopyHeader(src http.Header,dest *http.Header){
	for k,v:=range src{
		dest.Set(k,v[0])
	}
}


func ProxyUrl(w http.ResponseWriter,r *http.Request,url string){
	request, _ := http.NewRequest(r.Method, url, r.Body)
	CopyHeader(r.Header,&request.Header)
	request.Header.Add("x-forwarded-for",r.RemoteAddr)
	response, e := http.DefaultClient.Do(request)
	if response==nil{
		log.Println(e)
		return
	}
	dest := w.Header()
	CopyHeader(response.Header,&dest)
	w.WriteHeader(response.StatusCode)
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)
	w.Write(bytes)
}