package test

import (
	"encoding/base64"
	"github.com/go-ini/ini"
	"regexp"
	"testing"
)


func Test_Enc(t *testing.T){
 	str:="shuchang:123"
	s := base64.StdEncoding.EncodeToString([]byte(str))
	t.Log(s)
	//c2h1Y2hhbmc6MTIz
}

func Test_ini(t *testing.T){
	file, e := ini.Load("/Users/shuchang/Documents/GoPath/src/rev-proxy/env.ini")
	if e!=nil{
		t.Log(e)
		return
	}
	section, _ := file.GetSection("proxy")
	for _, sec := range section.ChildSections() {
		t.Log(sec.Name())
	}
}


func Test_reg(t *testing.T){
	str:="http://localhost:8080 weight=50"
	compile  := regexp.MustCompile(`(?<=(weight=))\d+`)
	matchString := compile.FindStringSubmatch(str)
	t.Log(matchString)
}