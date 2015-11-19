/*******************************************************************
 *  Copyright(c) 2000-2015 rjcb99
 *  All rights reserved.
 *
 *  文件名称: msgBox.go
 *  简要描述: 信箱
 *
 *  创建日期: 2015-11-18
 *  作者: ChenBo
 *  说明:
 *
 *  修改日期: 2015-11-19
 *  作者: ChenBo
 *  说明:
 ******************************************************************/
package main

import (
	"encoding/binary"
	"fmt"
	"game/utils"
	"net"
	"sync"
	"time"
)

//信箱
type MsgBox struct {
	BoxId     int
	MB        *utils.SafeQueue
	ExistMsg  chan bool
	MBM       *PostMan
	Conns     *utils.SafeStrMap
	ConnIndex int
	lock      *sync.RWMutex
}

//创建一个信箱，每个信箱限制2000条信息上限，超过上限，删除比较老的信息
func MakeNewMsgBox(id int, pm *PostMan) *MsgBox {
	mb := &MsgBox{
		BoxId:     id,
		MB:        utils.MakeNewSafeQueue(2000),
		ExistMsg:  make(chan bool),
		MBM:       pm,
		Conns:     utils.MakeNewSafeStrMap(),
		lock:      new(sync.RWMutex),
		ConnIndex: 0,
	}
	mb.Load()
	go mb.Sender()
	return mb
}

//关闭信箱
func (this *MsgBox) Close() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Conns.EachItem(func(key string, e interface{}) {
		if c, ok := e.(net.Conn); ok {
			c.Close()
		}
	})
	this.Conns.Clear()
	this.Save()
}

//存信息
func (this *MsgBox) AddMsg(msg *SimpleMsg) {
	if int(msg.MsgReceiver) == this.BoxId {
		this.MB.Push(msg)
		if this.Conns.Size() > 0 {
			this.ExistMsg <- true
		}
	}
}

//提供用户临时id
func (this *MsgBox) GetConnIndex() string {
	this.ConnIndex += 1
	if this.ConnIndex >= 1000000000 {
		this.ConnIndex = 1
	}
	return utils.I2s(this.ConnIndex)
}

//信箱托管用户，同一个信箱允许多个用户同时发送接收信息，并且按照先后顺序给用户编号1～999999999
func (this *MsgBox) AddConn(conn net.Conn) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if conn != nil {
		key := this.GetConnIndex()
		this.Conns.Set(key, conn)
		go this.Server(key, conn)
		if this.MB.Size() > 0 {
			this.ExistMsg <- true
		}
	}
}

//信箱托管的业务 接收用户发送的信息进行业务处理
func (this *MsgBox) Server(key string, conn net.Conn) {
	go this.Reader(key, conn)
}

//将信息转发给在线用户 没有在线用户就保存不发，有多个在线用户就同时发送
func (this *MsgBox) Sender() {
	for {
		_, ok := <-this.ExistMsg
		if !ok { //信箱关闭
			this.Close()
			return
		}
		for { //轮询
			if this.MB.Size() > 0 && this.Conns.Size() > 0 {
				flag := false
				sm := this.MB.Pop()
				if sm, ok := sm.(*SimpleMsg); ok && sm != nil {
					this.Conns.EachItem(func(key string, e interface{}) {
						if conn, ok := e.(net.Conn); ok && conn != nil {
							flag = true
							go this.SendMsg(key, conn, sm) //异步发送数据
						}
					})
				}
				if !flag {
					this.MB.PushFront(sm) //信息归还
				}
			} else {
				break
			}
		}
	}
}

//发送数据 3秒超时删除连接
func (this *MsgBox) SendMsg(key string, conn net.Conn, msg *SimpleMsg) {
	if conn != nil {
		conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
		_, err := conn.Write(msg.ToData()) //包续传先不管
		if err != nil {
			conn.Close()
			this.Conns.Remove(key)
			SysLog.PutLineAsLog(fmt.Sprintf("Msgbox%s SendMsg: %s ERROR:%s", this.BoxId, msg.MsgBody, err.Error()))
		}
	} else {
		this.Conns.Remove(key)
	}
}

//信箱的消息读取业务
func (this *MsgBox) Reader(key string, conn net.Conn) {
	if conn == nil {
		return
	}
	var buffer [MaxMsgLen * 10]byte
	bufferLen := 0
	for {
		conn.SetReadDeadline(time.Now().Add(time.Hour * 24))
		//接收数据
		n, err := conn.Read(buffer[bufferLen:])
		if err != nil {
			conn.Close()
			this.Conns.Remove(key)
			return
		}
		if err == nil && n > 0 {
			bufferLen += n
		}
		//处理数据
		for {
			if bufferLen < 4 {
				break
			}
			msglen := int(binary.LittleEndian.Uint32(buffer[0:4]))
			if msglen > MaxMsgLen {
				conn.Close()
				this.Conns.Remove(key)
				return
			} else {
				if bufferLen >= msglen {
					if sm := MakeNewSimpleMsg(); sm != nil {
						sm.FromBytes(buffer[0:msglen])
						copy(buffer[0:], buffer[msglen:])
						bufferLen -= msglen
						this.MBM.Sorting(sm)
					} else {
						conn.Close()
						this.Conns.Remove(key)
						return
					}
				} else {
					break
				}
			}
		}
	}
}

//落地数据重载 根据需求实现
func (this *MsgBox) Load() {
}

//数据落地 根据需求实现
func (this *MsgBox) Save() {
}

//数据过滤 一般是用正则式过滤敏感词
func (this *MsgBox) Filter() {
}
