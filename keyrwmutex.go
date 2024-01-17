package keyrwmutex

import (
	"hash/fnv"
	"runtime"
	"sync"

	"k8s.io/utils/keymutex"
)

var _ keymutex.KeyMutex = (*keyRWMutex)(nil)

type keyRWMutex struct {
	mutexes []sync.RWMutex
}

func New(n int) *keyRWMutex {
	if n <= 0 {
		n = runtime.NumCPU()
	}
	return &keyRWMutex{
		mutexes: make([]sync.RWMutex, n),
	}
}

func (km *keyRWMutex) LockKey(key string) {
	km.mutexes[km.hash(key)%uint32(len(km.mutexes))].Lock()
}

func (km *keyRWMutex) UnlockKey(key string) error {
	km.mutexes[km.hash(key)%uint32(len(km.mutexes))].Unlock()
	return nil
}

func (km *keyRWMutex) RLockKey(key string) {
	km.mutexes[km.hash(key)%uint32(len(km.mutexes))].RLock()
}

func (km *keyRWMutex) RUnlockKey(key string) error {
	km.mutexes[km.hash(key)%uint32(len(km.mutexes))].RUnlock()
	return nil
}

func (km *keyRWMutex) hash(id string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(id))
	return h.Sum32()
}
