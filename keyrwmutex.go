package keyrwmutex

import (
	"hash/fnv"
	"runtime"
	"sync"

	"k8s.io/utils/keymutex"
)

var _ keymutex.KeyMutex = (*KeyRWMutex)(nil)

type KeyRWMutex struct {
	mutexes []sync.RWMutex
}

func New(n int) *KeyRWMutex {
	if n <= 0 {
		n = runtime.NumCPU()
	}
	return &KeyRWMutex{
		mutexes: make([]sync.RWMutex, n),
	}
}

func (km *KeyRWMutex) LockKey(key string) {
	km.mutexes[km.hash(key)%uint32(len(km.mutexes))].Lock()
}

func (km *KeyRWMutex) UnlockKey(key string) error {
	km.mutexes[km.hash(key)%uint32(len(km.mutexes))].Unlock()
	return nil
}

func (km *KeyRWMutex) RLockKey(key string) {
	km.mutexes[km.hash(key)%uint32(len(km.mutexes))].RLock()
}

func (km *KeyRWMutex) RUnlockKey(key string) error {
	km.mutexes[km.hash(key)%uint32(len(km.mutexes))].RUnlock()
	return nil
}

func (km *KeyRWMutex) hash(id string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(id))
	return h.Sum32()
}
