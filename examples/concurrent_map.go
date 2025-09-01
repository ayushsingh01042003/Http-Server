package examples

import "sync"

type ConcurrentMap struct {
	m  map[int]int
	mu sync.RWMutex
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		m: make(map[int]int),
	}
}

func (cm *ConcurrentMap) Set(key int, val int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m[key] = val
}

func (cm *ConcurrentMap) Get(key int) (int, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	val, ok := cm.m[key]
	return val, ok
}

func (cm *ConcurrentMap) Delete(key int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.m, key)
}

