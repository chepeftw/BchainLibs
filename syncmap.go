package bchainlibs

import "sync"

type MapBlocks struct {
	sync.RWMutex
	m map[string]Packet
}

func (mapb MapBlocks) Get( key string ) Packet {
	mapb.RLock()
	n := mapb.m[key]
	mapb.RUnlock()
	return n
}

func (mapb MapBlocks) Has( key string ) bool {
	mapb.RLock()
	_, ok := mapb.m[ key ]
	mapb.RUnlock()
	return ok
}

func (mapb MapBlocks) Add( key string, packet Packet ) {
	mapb.Lock()
	mapb.m[ key ] = packet
	mapb.Unlock()
}

func (mapb MapBlocks) Del( key string ) {
	mapb.Lock()
	delete(mapb.m, key)
	mapb.Unlock()
}