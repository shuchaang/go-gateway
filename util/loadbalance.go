package util

import (
	"hash/crc32"
	"log"
	"math/rand"
	"time"
)

type LoadBalance struct {
	Servers []*HttpServer
}


var serverIndices []int

var RR_INDEX=0

const RAND  = "RAND"
const HASH  = "IP_HASH"
const WEIGHT = "WEIGHT"




func (this *LoadBalance) SelectAlg(alg string, servers []*HttpServer,reqIp string) *HttpServer {
	if this.Servers==nil||len(this.Servers)==0{
		log.Println("ALG:",alg)
		this.Servers=servers
	}
	switch alg {
	case RAND:
		return this.selectForRand()
	case HASH:
		return this.selectForHash(reqIp)
	case WEIGHT:
		return this.selectForWeight()
	default:
		return this.selectForRand()
	}
}

func (this *LoadBalance) selectForRand()*HttpServer{
	rand.Seed(time.Now().UnixNano())
	return this.Servers[rand.Intn(len(this.Servers))]
}

func (this *LoadBalance) selectForHash(ip string) *HttpServer{
	index:=int(crc32.ChecksumIEEE([]byte(ip)))%len(this.Servers)
	return this.Servers[index]
}

func (this *LoadBalance) selectForWeight() *HttpServer {
	if serverIndices==nil||len(serverIndices)==0{
		for index, item := range this.Servers {
			if item.Weight>0{
				for i:=0;i<item.Weight;i++{
					serverIndices=append(serverIndices,index)
				}
			}
		}
		log.Println("init rand weight:",serverIndices)
	}
	rand.Seed(time.Now().UnixNano())
	return this.Servers[rand.Intn(len(this.Servers))]
}

func (this *LoadBalance) selectForRR() *HttpServer {
	RR_INDEX=(RR_INDEX+1)%len(this.Servers)
	return this.Servers[RR_INDEX]
}