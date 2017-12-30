package riff

import "sort"

type Services map[string]*Service

type Service struct {
	Name    string
	Address string
	Version uint64
}

func (ss *Services) sort() []string {
	var keys = make([]string, 0, 0)
	for key, _ := range *ss {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}
