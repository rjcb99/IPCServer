package utils

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

//信箱For Client
type MsgBox struct {
	BoxId     int
	Ip        string
	Port      int
	MB        *SafeQueue
	Conn      net.Conn
	lock      *sync.RWMutex
	ConnIndex int
}

//创建一个信箱，每个信箱限制2000条信息上限，超过上限，删除比较老的信息
func MakeNewMsgBox(id int, ip string, port int) *MsgBox {
	mb := &MsgBox{
		BoxId:     id,
		Ip:        ip,
		Port:      port,
		MB:        MakeNewSafeQueue(2000),
		lock:      new(sync.RWMutex),
		ConnIndex: 0,
	}
	return mb
}

//关闭信箱
func (this *MsgBox) Close() {
	this.Save()
	this.Conn.Close()
	this.Conn = nil
	SysLog.PutLineAsLog(fmt.Sprintf("MsgBox Close! ClientNo:%d Begin to ReConn[%s:%d]...", this.BoxId, this.Ip, this.Port))
	//不死重连 5s 一次 这个可以有
	for {
		if !this.MakeConn() {
			SysLog.PutLineAsLog(fmt.Sprintf("MsgBox ReConn False! ClientNo:%d After 5s to ReConn[%s:%d]...", this.BoxId, this.Ip, this.Port))
			time.Sleep(time.Second * 5)
		}
	}
}

//信箱存信息
func (this *MsgBox) PushMsg(msg *SimpleMsg) {
	if int(msg.MsgReceiver) == this.BoxId {
		this.MB.Push(msg)
	}
}

//信箱取消息
func (this *MsgBox) PopMsg() *SimpleMsg {
	if sm, ok := this.MB.Pop().(*SimpleMsg); ok && sm != nil {
		return sm
	} else {
		return nil
	}
}

//连接服务
func (this *MsgBox) MakeConn() bool {
	//地址解析
	connStr := fmt.Sprintf("%s:%d", this.Ip, this.Port)
	addr, err := net.ResolveTCPAddr("tcp", connStr)
	if err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf(" ClientHandle error resolving tcp address: %s, %s", connStr, err.Error()))
		return false
	}
	//连接
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf(" ClientHandle error connecting to server@%s: %s", connStr, err.Error()))
		return false
	}
	this.Server(conn)
	this.SendMsg(0, "") //登陆用
	return true
}

//信箱托管的业务 接收用户发送的信息进行业务处理
func (this *MsgBox) Server(conn net.Conn) {
	this.Conn = conn
	go this.Reader(conn)
}

//发送数据 3秒超时删除连接
func (this *MsgBox) SendMsg(toId int, msg string) {
	if this.Conn != nil {
		go this.sendMsg(toId, msg)
	}
}

//发送数据 3秒超时删除连接
func (this *MsgBox) sendMsg(toId int, msg string) {
	if this.Conn != nil {
		this.lock.Lock()
		defer this.lock.Unlock()
		sm := MakeNewSimpleMsg().FromString(this.BoxId, toId, msg)
		SysLog.PutLineAsLog(fmt.Sprintf("Msgbox%d SendMsg: %s", this.BoxId, msg))
		this.Conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
		_, err := this.Conn.Write(sm.ToData()) //包续传先不管
		if err != nil {
			this.Conn.Close()
			SysLog.PutLineAsLog(fmt.Sprintf("Msgbox%d SendMsg: %s ERROR:%s", this.BoxId, msg, err.Error()))
		}
	}
}

//信箱的消息读取业务
func (this *MsgBox) Reader(conn net.Conn) {
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
			return
		}
		if err == nil && n > 0 {
			bufferLen += n
		}
		//处理数据 buffer
		for {
			if bufferLen < 4 {
				break
			}
			msglen := int(binary.LittleEndian.Uint32(buffer[0:4]))
			if msglen > MaxMsgLen { //包超长
				conn.Close()
				return
			} else {
				if bufferLen >= msglen {
					sm := MakeNewSimpleMsg()
					sm.FromBytes(buffer[0:msglen])
					copy(buffer[0:], buffer[msglen:])
					bufferLen -= msglen
					this.PushMsg(sm)
					SysLog.PutLineAsLog(fmt.Sprintf("[Recv From:%d]:%s", sm.MsgSender, sm.ToString()))
				} else {
					break
				}
			}
		}
	}
}

//落地数据重载 后续实现
func (this *MsgBox) Load() *MsgBox {
	return this
}

//数据落地 后续实现
func (this *MsgBox) Save() {
}

//数据过滤 暂无需求
func (this *MsgBox) Filter() {
}
