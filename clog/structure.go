package clog

import "sync"

//日志项
type LogItem struct {
	Term    int
	Command string
}

type Manager struct {
	logs []LogItem
	lock *sync.RWMutex
}

func (this *Manager) Init() {
	this.lock = new(sync.RWMutex)
}

func (this *Manager) Add(log LogItem) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.logs = append(this.logs, log)
	return len(this.logs) - 1
}

func (this *Manager) Extend(otherLogs []LogItem) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.logs = append(this.logs, otherLogs...)
	return len(this.logs) - 1
}

func (this *Manager) Size() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.logs)
}

func (this *Manager) Get(idx int) LogItem {
	size := this.Size()
	if size == 0 || idx >= size {
		return LogItem{}
	}
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.logs[idx]
}

func (this *Manager) GetFrom(idx int) []LogItem {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.logs[idx:]
}

func (this *Manager) RemoveFrom(idx int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.logs = this.logs[:idx]
}