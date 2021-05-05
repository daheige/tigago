// Package chanlock 基于chan实现trylock乐观锁
package chanlock

import (
	"time"
)

// DefaultLockTimeout 默认加锁超时时间20ms就认为加锁失败
var DefaultLockTimeout = 10 * time.Millisecond

// ChanLock chan lock
type ChanLock struct {
	ch chan struct{} // 空结构体
}

// NewChanLock 实例化一个通道空结构体锁对象
func NewChanLock() *ChanLock {
	return &ChanLock{
		ch: make(chan struct{}, 1), // 有缓冲通道
	}
}

// Lock 通道加锁,如果无法放入ch，该方法就会阻塞，直到ch通道锁释放为止
func (l *ChanLock) Lock() {
	l.ch <- struct{}{} // 这里是一个空结构体
}

// Unlock实现通道解锁
func (l *ChanLock) Unlock() {
	<-l.ch
}

// TryLock 乐观锁实现
func (l *ChanLock) TryLock(timeout ...time.Duration) bool {
	if len(timeout) > 0 && timeout[0] > 0 {
		return l.tryLockTimeout(timeout[0])
	}

	return l.tryLockTimeout(DefaultLockTimeout)
}

// TryLockTimeout 指定时间内的乐观锁
func (l *ChanLock) tryLockTimeout(timeout time.Duration) bool {
	if timeout == 0 {
		timeout = DefaultLockTimeout
	}

	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case l.ch <- struct{}{}:
			return true
		case <-ticker.C:
			return false
		}
	}
}
