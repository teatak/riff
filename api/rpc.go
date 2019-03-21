package api

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type RpcClient struct {
	rpc string
}

func (this *RpcClient) Services(name string, state StateType) (service Service) {
	client, err := NewClient(this.rpc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	err = client.Call("Query.Service", ParamService{Name: name, State: state}, &service)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func reserveAddress(url string) (string,string)  {
	prefix := "http://"
	serviceName := ""
	urls := strings.SplitN(url, "//", 2)
	if len(urls) == 0 {
		return "",url
	} else {
		prefix = urls[0] + "//"
		serviceName = urls[1]
		if prefix == "rpc://" {
			prefix = ""
		}
		if prefix == "tcp://" {
			prefix = ""
		}
		return prefix,serviceName
	}
}
/*
robin
url: http://serviceName or rpc://serviceName
http url return http://ip:port
rpc url only return ip:port
*/
func (this *RpcClient) Robin(url string) (string, error) {
	prefix,serviceName := reserveAddress(url)
	service := this.Services(serviceName, StateAlive)
	count := len(service.NestNodes)
	if count > 0 {
		r := GenerateNumber(0, count-1)
		return prefix + service.NestNodes[r].IP + ":" + strconv.Itoa(service.NestNodes[r].Port), nil
	}
	return "", errors.New("404")
}

var counter = make(map[string]int)

/*
round
url: http://serviceName or rpc://serviceName
http url return http://ip:port
rpc url only return ip:port
*/
func (this *RpcClient) Round(url string) (string, error) {
	prefix,serviceName := reserveAddress(url)
	service := this.Services(serviceName, StateAlive)
	count := len(service.NestNodes)
	if count > 0 {
		r := counter[serviceName]
		if r >= count {
			r = 0
		}
		u := prefix + service.NestNodes[r].IP + ":" + strconv.Itoa(service.NestNodes[r].Port)
		if r == count-1 {
			r = 0
		} else {
			r++
		}
		counter[serviceName] = r
		return u, nil
	}
	return "", errors.New("404")
}

func GenerateNumber(min, max int) int {
	i := rand.Intn(max-min) + min
	return i
}