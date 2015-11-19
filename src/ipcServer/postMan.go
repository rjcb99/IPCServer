/*******************************************************************
 *  Copyright(c) 2000-2015 rjcb99
 *  All rights reserved.
 *
 *  文件名称: postMan.go
 *  简要描述: 邮递员负责消息的转储
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
	"fmt"
	"game/utils"
	"net"
	"time"
)

//邮递员
type PostMan struct {
	MBS *utils.SafeStrMap //安全的信箱Map
}

func MakeNewPostMan() *PostMan {
	pm := &PostMan{
		MBS: utils.MakeNewSafeStrMap(),
	}
	pm.InitMsgBox()
	return pm
}

//初始化信箱 默认初始化预编号1～100的信箱 如果信箱落地 则加载落地内容
func (this *PostMan) InitMsgBox() {
	for i := 1; i <= 100; i++ {
		this.MBS.Set(utils.I2s(i), MakeNewMsgBox(i, this))
	}
}

//移除信箱 出于性能考虑 一般不这么做
func (this *PostMan) RemoveMsgBox(boxid string) {
	if MB := this.MBS.Get(boxid); MB != nil {
		this.MBS.Remove(boxid)
		if mb, ok := MB.(*MsgBox); ok && mb != nil {
			mb.Close()
		}
	}

}

/*
 * 消息分发
 * 特殊消息1：发给 0的消息为广播消息，将会发送给每一个连接用户，除了自己
 * 特殊消息2：发给 99999999一个用户id 返回该id的NAT地址 可用于P2P
 */
func (this *PostMan) Sorting(msg *SimpleMsg) {
	fromId := utils.I2s(int(msg.MsgSender))
	toId := utils.I2s(int(msg.MsgReceiver))
	if toId != "99999999" && toId != "0" && this.MBS.IsExist(toId) {
		if mb, ok := this.MBS.Get(toId).(*MsgBox); ok && mb != nil {
			go mb.AddMsg(msg)
		} else {
			SysLog.PutLineAsLog(fmt.Sprintf("[PostMan] Sorting MBS.Get(%s).(*MsgBox) Error.msg:%v", toId, msg))
		}
	} else {
		if msg.MsgReceiver > 0 && msg.MsgReceiver < 99999999 { //私人信息
			this.MBS.Set(toId, MakeNewMsgBox(int(msg.MsgReceiver), this))
			this.Sorting(msg)
		} else if msg.MsgReceiver == 0 { //群发消息
			this.MBS.EachItem(func(key string, e interface{}) {
				if key != fromId { //自己不能发给自己
					if mb, ok := e.(*MsgBox); ok && mb != nil {
						go mb.AddMsg(msg)
					} else {
						SysLog.PutLineAsLog(fmt.Sprintf("PostMan Sorting e.(*MsgBox) Error.msg:%v", fromId, msg))
					}
				}
			})
		} else if msg.MsgReceiver == 99999999 { //特殊消息，获取某个MB的远端地址
			if mb, ok := this.MBS.Get(fromId).(*MsgBox); ok && mb != nil {
				addr := this.GetRemoteAddr(string(msg.MsgBody[12:]))
				addr = string(msg.MsgBody[12:]) + ":" + addr //format[id:255.255.255.255:123456]
				msg.FromString(99999999, int(msg.MsgSender), addr)
				go mb.AddMsg(msg)
			} else {
				SysLog.PutLineAsLog(fmt.Sprintf("PostMan Sorting MBS.Get(%s).(*MsgBox) Error.msg:%v", fromId, msg))
			}
		} else { //消息扩展......
		}
	}
}

//得到用户的addr
func (this *PostMan) GetRemoteAddr(uid string) string {
	addr := ""
	mb := this.MBS.Get(uid)
	if mb != nil {
		if m, ok := mb.(*MsgBox); ok && m != nil {
			m.Conns.EachItemBreak(func(key string, e interface{}) bool {
				if e != nil {
					if conn, ok := e.(net.Conn); ok && conn != nil {
						addr = conn.RemoteAddr().String()
						return true
					}
				}
				return false
			})
		} else {
			SysLog.PutLineAsLog(fmt.Sprintf("PostMan getRemoteAddr mb.(*MsgBox) Error.uid:%d", uid))
		}
	}
	return addr
}

//客户分发，登陆后3秒内注册id，否则自动断开连接，id注册后连接转交给信箱托管
func (this *PostMan) AddConn(conn net.Conn) {
	buf := make([]byte, 12)
	SysLog.PutLineAsLog("用户[*]连接")
	conn.SetReadDeadline(time.Now().Add(time.Second * 3)) //登陆的时候只给3秒的机会,根据带宽自行调整，越小越安全
	n, err := conn.Read(buf)
	if err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("PostMan AddConn Read data Error: %s", err.Error()))
		conn.Close()
		return
	}
	//注册id
	if n != 12 { //长度不够就断开
		SysLog.PutLineAsLog(fmt.Sprintf("PostMan AddConn Read data Error:注册数据长度不够!"))
		conn.Close()
	} else { //长度有多的要保留
		sm := MakeNewSimpleMsg().FromBytes(buf[0:12])
		if sm.MsgSize == 12 && sm.MsgSender > 0 && sm.MsgReceiver == 0 {
			SysLog.PutLineAsLog(fmt.Sprintf("User:%d Reg!", sm.MsgSender))
			//查找信箱
			if this.MBS.IsExist(utils.I2s(int(sm.MsgSender))) { //存在信箱 Conn交给信箱托管
				if mb, ok := this.MBS.Get(utils.I2s(int(sm.MsgSender))).(*MsgBox); ok && mb != nil {
					go mb.AddConn(conn)
				} else {
					SysLog.PutLineAsLog(fmt.Sprintf("PostMan AddConn Get MsgBox Error!"))
					conn.Close()
				}
			} else { //不存在信箱 创建信箱 Conn交给信箱托管
				if mb := MakeNewMsgBox(int(sm.MsgSender), this); mb != nil {
					this.MBS.Set(utils.I2s(int(sm.MsgSender)), mb)
					go mb.AddConn(conn)
				} else {
					SysLog.PutLineAsLog(fmt.Sprintf("PostMan AddConn MakeNewMsgBox Error!"))
					conn.Close()
				}
			}
		} else {
			SysLog.PutLineAsLog(fmt.Sprintf("PostMan AddConn read data Error: 身份格式不对[%v]", sm))
			conn.Close()
		}
	}
}
