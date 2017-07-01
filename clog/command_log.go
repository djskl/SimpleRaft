package clog

import (
	"sync"
	"SimpleRaft/utils"
	"strings"
	"fmt"
)

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
	return len(this.logs)
}

func (this *Manager) Extend(preIdx int, otherLogs []LogItem) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	if otherLogs != nil && len(otherLogs) > 0{
		if preIdx == len(this.logs) {
			this.logs = append(this.logs, otherLogs...)
		}

		if preIdx < len(this.logs) {
			this.logs = this.logs[:preIdx]
			this.logs = append(this.logs, otherLogs...)
		}
	}
	return len(this.logs)
}

func (this *Manager) Size() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return len(this.logs)
}

func (this *Manager) Get(idx int) (LogItem, error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	size := len(this.logs)
	if idx < 1 || size == 0 || idx > size {
		return LogItem{}, utils.NotExistsError{"日志不存在"}
	}
	return this.logs[idx-1], nil
}

func (this *Manager) GetMany(beg int, end int) []LogItem {
	this.lock.RLock()
	defer this.lock.RUnlock()

	size := len(this.logs)
	if beg < 1 || beg > size {
		return nil
	}

	if end < beg || end > size + 1 {
		return nil
	}

	return this.logs[beg-1:end-1]

}

func (this *Manager) GetFrom(idx int) []LogItem {
	this.lock.RLock()
	defer this.lock.RUnlock()
	size := len(this.logs)
	if idx < 1 || idx > size {
		return nil
	}
	return this.logs[idx-1:]
}

func (this *Manager) RemoveFrom(idx int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.logs = this.logs[:idx]
}

func (this *Manager) ToString() string {
	this.lock.Lock()
	defer this.lock.Unlock()
	var rst []string
	for _, val := range this.logs{
		rst = append(rst, val.Command)
	}
	s := fmt.Sprintf("[%s]", strings.Join(rst, ", "))
	return s
}
