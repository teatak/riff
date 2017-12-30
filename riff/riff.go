package riff

import (
	"time"
	"math/rand"
)

//todo 今晚任务 1，序列化问题，json序列化支持

type Riff struct {
	Services []Service
	Nodes    []Node
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Create() (*Riff, error) {
	return nil, nil
}
