package bchainlibs

import (
	"sync"
	"time"
)

type MapBlocks struct {
	M map[string]Packet
	ST map[string]int64
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
	mapb.Lock()
	mapb.M[ key ] = packet
	mapb.ST[ key ] = time.Now().UnixNano()
	mapb.Unlock()
}

func (mapb MapBlocks) Del( key string ) (int64) {
	start := int64(0)
	mapb.Lock()
	start = mapb.ST[key]
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