package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func main(){

	c:=make(chan os.Signal)

	go func() {
		http.ListenAndServe(":9091", &webHandler1{})
	}()

	go func() {
		http.ListenAndServe(":9092",&webHandler2{})
	}()
	go func() {
		http.ListenAndServe(":9093", &webHandler2{})
	}()
	go func() {
		http.ListenAndServe(":9094", &webHandler2{})
	}()
	signal.Notify(c,os.Interrupt)
	s:=<-c
	log.Print(s)
}


type webHandler1 struct {

}



func (w *webHandler1) ServeHTTP(writer http.ResponseWriter,req *http.Request) {
	auth := req.Header.Get("Authorization")
	if auth==""{
		writer.Header().Set("WWW-Authenticate",`Basic realm='输入用户名密码'`)
		writer.WriteHeader(401)
		return
	}
	split := strings.Split(auth, " ")

	if len(split)==2&&split[0]=="Basic"{
		bytes, e := base64.StdEncoding.DecodeString(split[1])
		code := string(bytes)
		if e==nil&&code=="shuchang:123"{
			writer.Write([]byte("<h1>Login success</h1>"))
			return
		}
	}
	writer.Header().Set("WWW-Authenticate",`Basic realm='输入用户名密码'`)
	writer.WriteHeader(401)
	writer.Write([]byte("<h1>Login Failed</h1>"))
}

type webHandler2 struct {

}
func (wh *webHandler2)GetIp(reqest *http.Request)string{
	xff:=reqest.Header.Get("x-forwarded-for")
	if xff!=""{
		split := strings.Split(xff, ",")
		if len(split)>0&&split[0]!=""{
			return split[0]
		}
	}
	return reqest.RemoteAddr
}
func (w *webHandler2) ServeHTTP(writer http.ResponseWriter,http *http.Request) {
	writer.Write([]byte(fmt.Sprintf("<h1>%s</h1>",w.GetIp(http))))
}