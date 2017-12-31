package riff

import (
	"fmt"
	"testing"
)

func TestRiff(t *testing.T) {
	n1 := &Node{
		Name:       "n1",
		DataCenter: "dc1",
	}
	n2 := &Node{
		Name:       "n2",
		DataCenter: "dc1",
	}
	n3 := &Node{
		Name:       "n3",
		DataCenter: "dc1",
	}
	n4 := &Node{
		Name:       "n4",
		DataCenter: "dc1",
	}

	riff, _ := Create()

	s1 := &Service{Name: "s1", Address: "192.168.1.1:8080", Version: 122}
	s2 := &Service{Name: "s2", Address: "192.168.1.2:8081", Version: 12}
	s3 := &Service{Name: "s3", Address: "192.168.1.2:8082", Version: 44}
	s4 := &Service{Name: "s4", Address: "192.168.1.2:8083", Version: 55}

	riff.Link(n1, s1)
	riff.Link(n1, s2)
	riff.Link(n1, s3)
	riff.Link(n2, s1)

	s1 = &Service{Name: "s1", Address: "192.168.1.2:8080", Version: 4545}

	riff.Link(n3, s1)
	riff.Link(n3, s2)

	riff.Link(n4, s3)
	riff.Link(n4, s4)

	riff.Shutter()

	if n,ok := riff.Nodes["n1"]; ok {
		fmt.Println(n)
	}

	fmt.Println(riff.Nodes)
	//riff.AddService()
}
