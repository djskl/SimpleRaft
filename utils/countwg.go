package utils

import "sync"

type CountWG struct {
	wg sync.WaitGroup
	num int
	wk sync.Mutex
}

func (this *CountWG) Add(delta int) {
	this.wk.Lock()
	defer this.wk.Unlock()
	this.num += delta
	this.wg.Add(delta)
}

func (this *CountWG) Done() {
	this.wk.Lock()
	defer this.wk.Unlock()
	if this.num == 0 {
		return
	}
	this.wg.Done()
	this.num -= 1
}

func (this *CountWG) Size() int {
	return this.num
}

func (this *CountWG) Wait() {
	this.wg.Wait()
}
