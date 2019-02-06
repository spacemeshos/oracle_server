package main

import (
	"github.com/spacemeshos/go-spacemesh/eligibility"
	"log"
	"sync"
)

type WorldSwitch struct {
	mtx    sync.Mutex
	worlds map[uint64]*eligibility.FixedRolacle

	cntMtx        sync.Mutex
	clientCounter map[uint64]int // world -> size
}

func NewWorldSwitch() *WorldSwitch {
	return &WorldSwitch{worlds: make(map[uint64]*eligibility.FixedRolacle), clientCounter: make(map[uint64]int)}
}

func (ws *WorldSwitch) Get(id uint64) *eligibility.FixedRolacle {
	var world *eligibility.FixedRolacle
	ws.mtx.Lock()
	world, ok := ws.worlds[id]
	if !ok {
		log.Println("Creating world %d", id)
		world = eligibility.New()
		ws.worlds[id] = world
	}
	ws.mtx.Unlock()
	return world
}

func (ws *WorldSwitch) Register(id uint64, isHonest bool, client string) {
	ws.Get(id).Register(isHonest, client)
	ws.cntMtx.Lock()
	ws.clientCounter[id]++
	ws.cntMtx.Unlock()
}

func (ws *WorldSwitch) Unregister(id uint64, isHonest bool, client string) {
	ws.Get(id).Register(isHonest, client)
	ws.cntMtx.Lock()
	ws.clientCounter[id]--
	c := ws.clientCounter[id]
	ws.cntMtx.Unlock()
	if c <= 0 {
		ws.remove(id) // todo : delete counter this ? the universe is infinite
	}
}

func (ws *WorldSwitch) remove(id uint64) {
	ws.mtx.Lock()
	delete(ws.worlds, id)
	ws.mtx.Unlock()
	log.Println("World removed : ", id)
}
