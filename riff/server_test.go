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

	s := Server{}
	for i := 0; i < 100; i++ {
		s.SetNode(&Node{
			Name: "node" + strconv.Itoa(i),
		})
	}
	s.Shutter()
	e := time.Now()
	fmt.Printf("time: %v,snap: %s\n", e.Sub(st), s.SnapShot)
}
