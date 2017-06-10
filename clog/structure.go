package clog

import "sync"

//日志项
type Item struct {
	Term    int
	Command string
}

type Manager struct {
	logs []Item
	lock *sync.RWMutex
}

func (this *Manager) Init() {
	this.lock = new(sync.RWMutex)
}

func (this *Manager) Add(log Item) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.logs = append(this.logs, log)
}

func (this *Manager) Size() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.logs)
}

func (this *Manager) Get(idx int) Item {
	size := this.Size()
	if idx >= size {
		return nil
	}
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.logs[idx]
}

func (this *Manager) GetFrom(idx int) []Item {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.logs[idx:]
}

func (this *Manager) RemoveFrom(idx int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.logs = this.logs[:idx]
}