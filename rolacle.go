package main

import (
	"log"
	"math/rand"
	"oracle_server/pb"
	"sync"
)

type WorldSwitch struct {
	mtx    sync.Mutex
	worlds map[uint64]*RolacleSwitch
}

func NewWorldSwitch() *WorldSwitch {
	return &WorldSwitch{worlds: make(map[uint64]*RolacleSwitch)}
}

func (ws *WorldSwitch) Get(id uint64) *RolacleSwitch {
	var world *RolacleSwitch
	ws.mtx.Lock()
	world, ok := ws.worlds[id]
	if !ok {
		log.Println("Creating world %d", id)
		world = NewRolacleSwitch()
		world.kill = func() {
			ws.mtx.Lock()
			delete(ws.worlds, id)
			ws.mtx.Unlock()
			log.Println("World removed : ", id)
		}
		ws.worlds[id] = world
	}
	ws.mtx.Unlock()
	return world
}

type RolacleSwitch struct {
	mtx     sync.Mutex
	clients map[string]struct{}

	kill func()

	instLock  sync.Mutex
	instances map[int64]map[string]struct{}
}

func NewRolacleSwitch() *RolacleSwitch {
	return &RolacleSwitch{clients: make(map[string]struct{}), instances: make(map[int64]map[string]struct{})}
}

func (rc *RolacleSwitch) Register(pubkey string) {
	rc.mtx.Lock()
	if _, exist := rc.clients[pubkey]; exist {
		rc.mtx.Unlock()
		return
	}

	rc.clients[pubkey] = struct{}{}
	rc.mtx.Unlock()
}

func (rc *RolacleSwitch) Unregister(pubkey string) {
	rc.mtx.Lock()
	delete(rc.clients, pubkey)
	l := len(rc.clients)
	rc.mtx.Unlock()
	if l == 0 {
		rc.kill()
	}
}

func (rc *RolacleSwitch) Validate(instanceID int64, committeeSize int, proof string) bool {
	rc.instLock.Lock()
	rolacle, ok := rc.instances[instanceID]
	if !ok {
		elmap := rc.createEligibilityMap(instanceID, committeeSize)
		rc.instances[instanceID] = elmap
		rc.instLock.Unlock()
		return rc.Validate(instanceID, committeeSize, proof)
	}
	rc.instLock.Unlock()
	_, valid := rolacle[proof]
	return valid
}

func (rc *RolacleSwitch) ValidateMap(instanceID int64, committeeSize int, proof string) *pb.ValidList {
	rc.instLock.Lock()
	rolacle, ok := rc.instances[instanceID]
	if ok {
		rc.instLock.Unlock()
		return MapToList(rolacle)
	}
	elmap := rc.createEligibilityMap(instanceID, committeeSize)
	rc.instances[instanceID] = elmap
	rc.instLock.Unlock()
	return MapToList(elmap)
}

func MapToList(elgmap map[string]struct{}) *pb.ValidList {
	vl := &pb.ValidList{}
	for k, _ := range elgmap {
		vl.IDs = append(vl.IDs, k)
	}
	return vl
}

// USE ONLY FROM `Validate`
func (rc *RolacleSwitch) createEligibilityMap(instanceID int64, committeeSize int) map[string]struct{} {
	clients := []string{}
	rc.mtx.Lock()
	l := len(rc.clients)

	if l < committeeSize {
		committeeSize = l
		// todo different codepath that just picks all registered as eligible
	}

	for k := range rc.clients {
		clients = append(clients, k)
	}
	rc.mtx.Unlock()

	selected := make(map[string]struct{})

	seed := rand.New(rand.NewSource(instanceID))

	for i := 0; i < committeeSize; i++ {
		randclient := clients[seed.Int31n(int32(len(clients)))]
		// ensure uniqueness
		_, ok := selected[randclient]
		for ok {
			randclient = clients[seed.Int31n(int32(len(clients)))]
			_, ok = selected[randclient]
		}

		selected[randclient] = struct{}{}
	}

	return selected
}
