package api

import (
	"errors"
	"strings"
)

type Client interface {
	Services(serviceName string, state StateType) Service
	Robin(url string) (string, error)
	Round(url string) (string, error)
	Hash(url, key string) (string, error)
	HashRing(url, key string) (string, error)
}

/*
url : riff://ip:port
return Client
*/
func RiffClient(url string) (Client, error) {
	if strings.Index(url, "riff://") == 0 {
		rpc := strings.Replace(url, "riff://", "", 1)
		return &RpcClient{rpc}, nil
	} else {
		return nil, errors.New("not support")
	}
}
