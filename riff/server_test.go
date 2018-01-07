package riff

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	//riff.AddService()
	st := time.Now()

	s := Server{
		Nodes: make(map[string]*Node),
	}
	for i := 0; i < 100; i++ {
		s.AddNode(&Node{
			Id:   strconv.Itoa(i),
			Name: "aaad",
		})
	}
	s.Shutter()
	e := time.Now()
	fmt.Printf("time: %v,snap: %s\n", e.Sub(st), s.SnapShot)
}
