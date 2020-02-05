package util

import (
	"github.com/go-ini/ini"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ProxyConfig struct {
	Mapping map[string][]*HttpServer
	LoadBalance string
}


type HttpServer struct {
	Host string
	Weight int
}

var Config *ProxyConfig

type EnvConfig *os.File

var regWeight =regexp.MustCompile(`weight=(\d+)`)


func init(){
	Config=&ProxyConfig{
		Mapping: make(map[string][]*HttpServer),
		LoadBalance:"rand",
	}
	EnvConfig, e := ini.Load("env.ini")
	if e!=nil{
		panic(e)
		return
	}

	main, _ := EnvConfig.GetSection("proxy")
	if main!=nil{
		sections := main.ChildSections()
		for _, sec := range sections {
			path, _ := sec.GetKey("path")
			us, _ := sec.GetKey("upstream")
			hosts:=make([]*HttpServer,0)
			if path!=nil&&us!=nil{
				split := strings.Split(us.Value(), ",")
				for _, host := range split {
					match := regWeight.FindStringSubmatch(host)
					weight:=0
					if len(match)==2{
						weight,_=strconv.Atoi(match[1])
					}
					hosts=append(hosts,&HttpServer{
						Host:host,
						Weight:weight,
					})
				}
			}
			Config.Mapping[path.Value()]=hosts
		}
	}
	lb, _ := EnvConfig.GetSection("loadBalance")
	v, _ := lb.GetKey("method")
	if v!=nil&&v.Value()!=""{
		Config.LoadBalance=strings.ToUpper(v.Value())
	}
	log.Print("配置文件加载:")
	log.Println(Config)
}
