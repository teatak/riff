package api

import (
	"crypto/md5"
	"fmt"
	"math"
	"sort"
)

type HashKey uint32
type HashKeyOrder []HashKey

func (h HashKeyOrder) Len() int           { return len(h) }
func (h HashKeyOrder) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h HashKeyOrder) Less(i, j int) bool { return h[i] < h[j] }

type consistentHash struct {
	ring       map[HashKey]string
	sortedKeys []HashKey
	nodes      []string
	weights    map[string]int
}

func newConsistentHash(nodes []string) *consistentHash {
	hash := &consistentHash{
		ring:       make(map[HashKey]string),
		sortedKeys: make([]HashKey, 0),
		nodes:      nodes,
		weights:    make(map[string]int),
	}
	hash.generateCircle()
	return hash
}

func newWithWeights(weights map[string]int) *consistentHash {
	nodes := make([]string, 0, len(weights))
	for node := range weights {
		nodes = append(nodes, node)
	}
	hash := &consistentHash{
		ring:       make(map[HashKey]string),
		sortedKeys: make([]HashKey, 0),
		nodes:      nodes,
		weights:    weights,
	}
	hash.generateCircle()
	return hash
}

func (h *consistentHash) Size() int {
	return len(h.nodes)
}

func (h *consistentHash) UpdateWithWeights(weights map[string]int) {
	nodesChgFlg := false
	if len(weights) != len(h.weights) {
		nodesChgFlg = true
	} else {
		for node, newWeight := range weights {
			oldWeight, ok := h.weights[node]
			if !ok || oldWeight != newWeight {
				nodesChgFlg = true
				break
			}
		}
	}

	if nodesChgFlg {
		newhring := newWithWeights(weights)
		h.weights = newhring.weights
		h.nodes = newhring.nodes
		h.ring = newhring.ring
		h.sortedKeys = newhring.sortedKeys
	}
}

func (h *consistentHash) generateCircle() {
	totalWeight := 0
	for _, node := range h.nodes {
		if weight, ok := h.weights[node]; ok {
			totalWeight += weight
		} else {
			totalWeight += 1
			h.weights[node] = 1
		}
	}

	for _, node := range h.nodes {
		weight := h.weights[node]

		factor := math.Floor(float64(40*len(h.nodes)*weight) / float64(totalWeight))

		for j := 0; j < int(factor); j++ {
			nodeKey := fmt.Sprintf("%s-%d", node, j)
			bKey := hashDigest(nodeKey)

			for i := 0; i < 3; i++ {
				key := hashVal(bKey[i*4 : i*4+4])
				h.ring[key] = node
				h.sortedKeys = append(h.sortedKeys, key)
			}
		}
	}

	sort.Sort(HashKeyOrder(h.sortedKeys))
}

func (h *consistentHash) GetNode(stringKey string) (node string, ok bool) {
	pos, ok := h.GetNodePos(stringKey)
	if !ok {
		return "", false
	}
	return h.ring[h.sortedKeys[pos]], true
}

func (h *consistentHash) GetNodePos(stringKey string) (pos int, ok bool) {
	if len(h.ring) == 0 {
		return 0, false
	}

	key := h.GenKey(stringKey)

	nodes := h.sortedKeys
	pos = sort.Search(len(nodes), func(i int) bool { return nodes[i] > key })

	if pos == len(nodes) {
		// Wrap the search, should return first node
		return 0, true
	} else {
		return pos, true
	}
}

func (h *consistentHash) GenKey(key string) HashKey {
	bKey := hashDigest(key)
	return hashVal(bKey[0:4])
}

func (h *consistentHash) GetNodes(stringKey string, size int) (nodes []string, ok bool) {
	pos, ok := h.GetNodePos(stringKey)
	if !ok {
		return nil, false
	}

	if size > len(h.nodes) {
		return nil, false
	}

	returnedValues := make(map[string]bool, size)
	//mergedSortedKeys := append(h.sortedKeys[pos:], h.sortedKeys[:pos]...)
	resultSlice := make([]string, 0, size)

	for i := pos; i < pos+len(h.sortedKeys); i++ {
		key := h.sortedKeys[i%len(h.sortedKeys)]
		val := h.ring[key]
		if !returnedValues[val] {
			returnedValues[val] = true
			resultSlice = append(resultSlice, val)
		}
		if len(returnedValues) == size {
			break
		}
	}

	return resultSlice, len(resultSlice) == size
}

func (h *consistentHash) AddNode(node string) *consistentHash {
	return h.AddWeightedNode(node, 1)
}

func (h *consistentHash) AddWeightedNode(node string, weight int) *consistentHash {
	if weight <= 0 {
		return h
	}

	if _, ok := h.weights[node]; ok {
		return h
	}

	nodes := make([]string, len(h.nodes), len(h.nodes)+1)
	copy(nodes, h.nodes)
	nodes = append(nodes, node)

	weights := make(map[string]int)
	for eNode, eWeight := range h.weights {
		weights[eNode] = eWeight
	}
	weights[node] = weight

	hashRing := &consistentHash{
		ring:       make(map[HashKey]string),
		sortedKeys: make([]HashKey, 0),
		nodes:      nodes,
		weights:    weights,
	}
	hashRing.generateCircle()
	return hashRing
}

func (h *consistentHash) UpdateWeightedNode(node string, weight int) *consistentHash {
	if weight <= 0 {
		return h
	}

	/* node is not need to update for node is not existed or weight is not changed */
	if oldWeight, ok := h.weights[node]; (!ok) || (ok && oldWeight == weight) {
		return h
	}

	nodes := make([]string, len(h.nodes))
	copy(nodes, h.nodes)

	weights := make(map[string]int)
	for eNode, eWeight := range h.weights {
		weights[eNode] = eWeight
	}
	weights[node] = weight

	hashRing := &consistentHash{
		ring:       make(map[HashKey]string),
		sortedKeys: make([]HashKey, 0),
		nodes:      nodes,
		weights:    weights,
	}
	hashRing.generateCircle()
	return hashRing
}
func (h *consistentHash) RemoveNode(node string) *consistentHash {
	/* if node isn't exist in hashring, don't refresh hashring */
	if _, ok := h.weights[node]; !ok {
		return h
	}

	nodes := make([]string, 0)
	for _, eNode := range h.nodes {
		if eNode != node {
			nodes = append(nodes, eNode)
		}
	}

	weights := make(map[string]int)
	for eNode, eWeight := range h.weights {
		if eNode != node {
			weights[eNode] = eWeight
		}
	}

	hashRing := &consistentHash{
		ring:       make(map[HashKey]string),
		sortedKeys: make([]HashKey, 0),
		nodes:      nodes,
		weights:    weights,
	}
	hashRing.generateCircle()
	return hashRing
}

func hashVal(bKey []byte) HashKey {
	return ((HashKey(bKey[3]) << 24) |
		(HashKey(bKey[2]) << 16) |
		(HashKey(bKey[1]) << 8) |
		(HashKey(bKey[0])))
}

func hashDigest(key string) [md5.Size]byte {
	return md5.Sum([]byte(key))
}
