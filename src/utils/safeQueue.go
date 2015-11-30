package utils

import (
	"container/list"
	"sync"
)

const (
	MaxNumDefalt = 1000
)

type SafeQueue struct {
	list   *list.List
	maxNum int
	lock   *sync.RWMutex
}

func MakeNewSafeQueue(maxnum int) *SafeQueue {
	sq := &SafeQueue{
		list:   list.New(),
		lock:   new(sync.RWMutex),
		maxNum: maxnum,
	}
	if maxnum <= 0 {
		sq.maxNum = MaxNumDefalt
	}
	return sq
}

//入队列 策略：损失旧的消息
func (this *SafeQueue) Push(e interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.list.PushBack(e)
	if this.list.Len() > this.maxNum {
		this.list.Remove(this.list.Front())
	}
}

//入队列到队首
func (this *SafeQueue) PushFront(e interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.list.PushFront(e)
}

//出队列
func (this *SafeQueue) Pop() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	e := this.list.Front()
	if e != nil {
		this.list.Remove(e)
		return e.Value
	}
	return nil
}

//得到队列最老元素
func (this *SafeQueue) Pick() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	e := this.list.Front()
	if e != nil {
		return e.Value
	}
	return nil
}

//删除某个元素
func (this *SafeQueue) Remove(ee interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for e := this.list.Front(); e != nil; e = e.Next() {
		if ee == e.Value {
			this.list.Remove(e)
			return
		}
	}
}

func (this *SafeQueue) Size() int {
	return this.list.Len()
}

//遍历
func (this *SafeQueue) EachItem(dealFun func(interface{})) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for e := this.list.Front(); e != nil; e = e.Next() {
		dealFun(e.Value)
	}
}

//遍历 可根据情况中断
func (this *SafeQueue) EachItemBreak(dealFun func(interface{}) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for e := this.list.Front(); e != nil; e = e.Next() {
		r := dealFun(e.Value)
		if r == true {
			break
		}
	}
}

//清空队列
func (this *SafeQueue) Clear() {
	this.lock.Lock()
	this.lock.Unlock()
	this.list.Init()
}
