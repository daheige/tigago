/**
* Package mutexlock mutex trylock.
* 在sync.Mutex基础上，实现乐观锁TryLock
* go1.18+后在sync Mutex已经实现了TryLock
 */
package mutexlock

import (
	"sync"
)

// NewMutexLock 创建lock实例
func NewMutexLock() *Mutex {
	return &Mutex{}
}

// Mutex mutex
type Mutex struct {
	in sync.Mutex
}

// Lock 加锁
func (m *Mutex) Lock() {
	m.in.Lock()
}

// Unlock 解锁
func (m *Mutex) Unlock() {
	m.in.Unlock()
}

// TryLock 尝试枷锁
func (m *Mutex) TryLock() bool {
	// const mutexLocked = 1 << iota
	// return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.in)), 0, mutexLocked)
	//
	// This method has been available since go1.18, so the official design is directly adopted.
	return m.in.TryLock()
}
