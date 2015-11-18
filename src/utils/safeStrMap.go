package utils

import (
	"sync"
)

//线程安全Map
type SafeStrMap struct {
	m    map[string]interface{}
	lock *sync.RWMutex
}

//构造
func MakeNewSafeStrMap() *SafeStrMap {
	return &SafeStrMap{
		m:    make(map[string]interface{}),
		lock: new(sync.RWMutex),
	}
}

//插入和更新一个元素
func (this *SafeStrMap) Set(key string, e interface{}) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if val, ok := this.m[key]; !ok {
		this.m[key] = e
	} else if val != e {
		this.m[key] = e
	} else {
		return false
	}
	return true
}

//删除一个元素
func (this *SafeStrMap) Remove(key string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.m, key)
	return true
}

//查找一个元素
func (this *SafeStrMap) Get(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	if val, ok := this.m[key]; ok {
		return val
	}
	return nil
}

//存在
func (this *SafeStrMap) IsExist(key string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.m[key]; !ok {
		return false
	}
	return true
}

//整个map长度
func (this *SafeStrMap) Size() int {
	this.lock.Lock()
	defer this.lock.Unlock()
	return len(this.m)
}

//清除整个map
func (this *SafeStrMap) Clear() {
	this.m = make(map[string]interface{})
}

//遍历 range对map而言实际上是无序的随机的
func (this *SafeStrMap) EachItem(dealFun func(string, interface{})) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for key, value := range this.m {
		dealFun(key, value)
	}
}

//遍历 可根据情况中断
func (this *SafeStrMap) EachItemBreak(dealFun func(string, interface{}) bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for key, value := range this.m {
		r := dealFun(key, value)
		if r == true {
			break
		}
	}
}

//模拟pop
func (this *SafeStrMap) Pop() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	key := ""
	for k, _ := range this.m {
		key = k
		break
	}
	if val, ok := this.m[key]; ok {
		delete(this.m, key)
		return val
	} else {
		return nil
	}
}

//模拟push
func (this *SafeStrMap) Push(key string, e interface{}) bool {
	return this.Set(key, e)
}
