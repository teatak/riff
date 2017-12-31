package riff

import (
	"fmt"
	"net"
	"testing"
)

func TestRiff(t *testing.T) {
	n1 := &Node{
		Name:       "n1",
		DataCenter: "dc1",
		IP:         net.ParseIP("192.168.1.1"),
		Port:       8530,
	}
	n2 := &Node{
		Name:       "n2",
		DataCenter: "dc1",
		IP:         net.ParseIP("192.168.1.2"),
		Port:       8530,
	}
	n3 := &Node{
		Name:       "n3",
		DataCenter: "dc1",
		IP:         net.ParseIP("192.168.1.3"),
		Port:       8530,
	}
	n4 := &Node{
		Name:       "n4",
		DataCenter: "dc1",
		IP:         net.ParseIP("192.168.1.4"),
		Port:       8530,
	}

	riff, _ := Create("node")

	s1 := &Service{Name: "s1", IP: net.ParseIP("192.168.1.1"), Port: 8080}
	s2 := &Service{Name: "s2", IP: net.ParseIP("192.168.1.2"), Port: 8080}
	s3 := &Service{Name: "s3", IP: net.ParseIP("192.168.1.3"), Port: 8080}
	s4 := &Service{Name: "s4", IP: net.ParseIP("192.168.1.4"), Port: 8080}

	riff.Link(n1, s1)
	riff.Link(n1, s2)
	riff.Link(n1, s3)
	riff.Link(n2, s1)

	riff.Link(n3, s1)
	riff.Link(n3, s2)

	riff.Link(n4, s3)
	riff.Link(n4, s4)

	riff.Shutter()

	fmt.Println(riff.Nodes)
	//riff.AddService()
}
