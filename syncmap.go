package bchainlibs

import (
	"sync"
)

type MapBlocks struct {
	M map[string]Packet
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
	mapb.Unlock()
}

func (mapb MapBlocks) Del( key string ) {
	mapb.Lock()
	delete(mapb.M, key)
	mapb.Unlock()
}

func (mapb MapBlocks) String() string {
	str := "{"

	mapb.RLock()
	length := len( mapb.M )
	i := 0
	for k, _ := range mapb.M {
		str += k
		i++
		if i < length { str += ", " }
	}
	mapb.RUnlock()

	str += "}"
	return str
}