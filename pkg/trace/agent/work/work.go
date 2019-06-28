package work

import (
	"sync"
	"sync/atomic"
)

type Pool struct {
	pool chan chan func()
	stop chan struct{}
	wg   sync.WaitGroup
	size int64
}

func NewPool(max int) *Pool {
	return &Pool{
		pool: make(chan chan func(), max),
		stop: make(chan struct{}),
	}
}

func (p *Pool) Run(fn func()) {
	select {
	case ch := <-p.pool:
		ch <- fn
	default:
		p.addWorker() <- fn
	}
}

func (p *Pool) Size() int64 {
	return atomic.LoadInt64(&p.size)
}

func (p *Pool) Stop() {
	close(p.stop)
	p.wg.Wait()
}

func (p *Pool) addWorker() chan func() {
	ch := make(chan func())
	p.wg.Add(1)
	atomic.AddInt64(&p.size, 1)
	go func() {
		defer atomic.AddInt64(&p.size, -1)
		defer p.wg.Done()
		for {
			select {
			case fn := <-ch:
				fn()
			case <-p.stop:
				return
			}
			select {
			case p.pool <- ch:
			default:
				// worker pool is full of idle workers
				return
			}
		}
	}()
	return ch
}
