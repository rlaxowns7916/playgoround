package problem

import (
	"context"
	"errors"
	"sync"
)

type Task func(context.Context) error

var (
	ErrClosed         = errors.New("pool closed")
	ErrSubmitCanceled = errors.New("submit canceled")
	ErrPoolCanceled   = errors.New("pool context canceled")
)

type Pool struct {
	ctx       context.Context
	wg        sync.WaitGroup
	taskQueue chan Task
	workQueue chan Task
	quit      chan struct{}
	once      sync.Once
}

func NewPool(ctx context.Context, workers int) *Pool {
	p := &Pool{
		ctx:       ctx,
		taskQueue: make(chan Task),
		workQueue: make(chan Task, workers),
		quit:      make(chan struct{}),
	}

	// boss
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(p.workQueue)
		for {
			select {
			case <-p.ctx.Done():
				return
			case <-p.quit:
				return
			case t := <-p.taskQueue:
				select {
				case <-p.ctx.Done():
					return
				case p.workQueue <- t:
				}
			}
		}
	}()

	//workQueue
	p.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-p.ctx.Done():
					return
				case w, ok := <-p.workQueue:
					if !ok {
						return
					}
					_ = w(p.ctx)
				}
			}
		}()
	}

	return p
}

func (p *Pool) Submit(ctx context.Context, t Task) error {
	select {
	case <-p.ctx.Done():
		return ErrPoolCanceled
	case <-p.quit:
		return ErrClosed
	case <-ctx.Done():
		return ErrSubmitCanceled
	case p.taskQueue <- t:
		return nil
	}
}

func (p *Pool) Close() {
	p.once.Do(func() { close(p.quit) })
}

func (p *Pool) Wait() error {
	p.wg.Wait()
	if p.ctx.Err() != nil {
		return ErrPoolCanceled
	}
	return nil
}
