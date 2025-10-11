package problem

import (
	"sync"
)

type Resource struct{}
type Lazy struct {
	once sync.Once
	res  *Resource
	err  error
	init func() (*Resource, error)
}

func NewLazy(f func() (*Resource, error)) *Lazy {
	return &Lazy{
		once: sync.Once{},
		init: f,
		res:  nil,
		err:  nil,
	}
}

func (l *Lazy) Get() (*Resource, error) {
	l.once.Do(func() {
		l.res, l.err = l.init()
	})

	return l.res, l.err
}
