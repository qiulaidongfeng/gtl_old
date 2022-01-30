// lmrwmutex
package gsync

import (
	"errors"
	"runtime"
	"sync/atomic"
)

var (
	Noreadlock  error = errors.New("No read lock")
	Nowritelock error = errors.New("No write lock")
)

const (
	nolock    = 0
	writelock = -1
)

type LMrwmutex struct {
	nm int64
}

func (m *LMrwmutex) Lock() {
	for {
		ok := atomic.CompareAndSwapInt64(&m.nm, nolock, writelock)
		if ok == true {
			break
		}
		runtime.Gosched()
	}
}

func (m *LMrwmutex) Unlock() {
	nm := atomic.LoadInt64(&m.nm)
	if nm >= 0 {
		panic(Nowritelock)
	}
	for {
		ok := atomic.CompareAndSwapInt64(&m.nm, writelock, nolock)
		if ok == true {
			break
		}
		runtime.Gosched()
	}
}

func (m *LMrwmutex) RLock() {
	for {
		ok := atomic.CompareAndSwapInt64(&m.nm, writelock, writelock)
		if ok == true {
			runtime.Gosched()
			continue
		}
		nm := atomic.AddInt64(&m.nm, 1)
		if nm > 0 {
			break
		}
	}
}

func (m *LMrwmutex) RUnlock() {
	nm := atomic.LoadInt64(&m.nm)
	if nm == -1 {
		panic(Noreadlock)
	} else if nm == 0 {
		panic(Noreadlock)
	}
	atomic.AddInt64(&m.nm, -1)
}
