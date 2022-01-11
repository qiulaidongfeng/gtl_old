package stack

import (
	"sync"
	"sync/atomic"
)

type slicestack struct {
	slice []interface{}
	size  *uint64
	mutex sync.RWMutex
}

func Newslicestack() slicestack {
	s := slicestack{
		slice: make([]interface{}, 0, 2),
		size:  new(uint64),
	}
	return s
}

func (s *slicestack) Push(x interface{}) error {
	s.slice = append(s.slice[:((*s.size)+1)], x)
	*s.size++
	return nil
}

func (s *slicestack) Tspush(x interface{}) error {
	s.mutex.Lock()
	s.slice = append(s.slice[:((*s.size)+1)], x)
	*s.size += 1
	s.mutex.Unlock()
	return nil
}

func (s *slicestack) Pop() (x interface{}, err error) {
	if (*s.size) == 0 {
		err = StackEmpty
		return x, err
	}
	x = s.slice[*s.size]
	*s.size -= 1
	return x, nil
}

func (s *slicestack) Tspop() (x interface{}, err error) {
	s.mutex.Lock()
	if *s.size == 0 {
		err = StackEmpty
		return x, err
	}
	x = s.slice[*s.size]
	*s.size -= 1
	s.mutex.Unlock()
	return x, nil

}

func (s *slicestack) Size() uint64 {
	return *s.size
}

func (s *slicestack) Tssize() uint64 {
	return atomic.LoadUint64(s.size)
}

func (s *slicestack) Clear() error {
	s.slice = make([]interface{}, 0, 2)
	*s.size = 0
	return nil
}

func (s *slicestack) Tsclear() error {
	s.mutex.Lock()
	s.slice = make([]interface{}, 0, 2)
	*s.size = 0
	s.mutex.Unlock()
	return nil
}

func (s slicestack) Look(size uint64) (interface{}, error) {
	if *s.size < size {
		return nil, StackSizeExceeded
	}
	return s.slice[size], nil
}

func (s slicestack) Tslook(size uint64) (interface{}, error) {
	s.mutex.Unlock()
	if *s.size < size {
		return nil, StackSizeExceeded
	}
	s.mutex.RUnlock()
	return s.slice[size], nil
}
