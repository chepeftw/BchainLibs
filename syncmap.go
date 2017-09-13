package bchainlibs

import (
	"sync"
	"time"
)

type MapBlocks struct {
	M map[string]Packet // Map itself
	ST map[string]int64	// Start Times
	HG map[string]int64	// Hashes Generated
	sync.RWMutex
}

func (mapb MapBlocks) Get( key string ) Packet {
	mapb.RLock()
	n := mapb.M[key]
	mapb.RUnlock()
	return n
}

func (mapb MapBlocks) Has( key string ) bool {
	mapb.RLock()
	_, ok := mapb.M[ key ]
	mapb.RUnlock()
	return ok
}

func (mapb MapBlocks) Add( key string, packet Packet ) {
	now := time.Now().UnixNano()
	mapb.Lock()
		mapb.M[ key ] = packet
		mapb.ST[ key ] = now
	mapb.Unlock()
}

func (mapb MapBlocks) Del( key string ) int64 {
	start := mapb.ST[key]
	mapb.Lock()
		delete(mapb.M, key)
		delete(mapb.ST, key)
	mapb.Unlock()

	return start
}

func (mapb MapBlocks) String() string {
	str := "{"

	mapb.RLock()
	length := len( mapb.M )
	i := 0
	for k := range mapb.M {
		str += k
		i++
		if i < length { str += ", " }
	}
	mapb.RUnlock()

	str += "}"
	return str
}

func (mapb MapBlocks) AddHashesCount( key string, count int64 ) {
	mapb.HG[ key ] += count
}

func (mapb MapBlocks) GetDelHashesCount( key string ) int64 {
	total := mapb.HG[ key ]
	delete(mapb.HG, key)
	return total
}