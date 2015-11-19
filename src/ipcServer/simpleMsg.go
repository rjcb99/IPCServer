/*******************************************************************
 *  Copyright(c) 2000-2015 rjcb99
 *  All rights reserved.
 *
 *  文件名称: simpleMsg.go
 *  简要描述: 简单[]byte消息
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
)

const (
	MaxMsgLen = 65536
)

//信息，一个简单的[]byte信息结构 限制最大长度65536
type SimpleMsg struct {
	MsgSize     uint32
	MsgSender   uint32
	MsgReceiver uint32 // 0:broadcast
	MsgBody     []byte //原始数据，包含头部，用于转发
}

/*
 *创建一条信息
 *1.[]byte格式：[4字节长度 + 4字节发送人id + 4字节接受人id + 正文] 其中长度 ＝ 12 + 正文字节长度
 *2.长度不能超过65536字节
 *3.socket连接成功3秒内要发身份验证消息，格式：12字节长度,正文为空，发送人id为验证id，接收人id可以任意
 *4.验证成功后发送消息根据接收人id分发，从1开始，可以自己给自己发;如果接收人id为0，表示群发该消息，但是不会发给自己。
 */
func MakeNewSimpleMsg() *SimpleMsg {
	return &SimpleMsg{
		MsgSize:     0,
		MsgSender:   0,
		MsgReceiver: 0,
		MsgBody:     []byte{},
	}
}

//[]byte -> SimpleMsg
func (this *SimpleMsg) FromBytes(buf []byte) *SimpleMsg {
	this.MsgSize = 0
	this.MsgSender = 0
	this.MsgReceiver = 0
	this.MsgBody = []byte{}
	if len(buf) < 12 {
		return this
	} else {
		this.MsgSize = binary.LittleEndian.Uint32(buf[0:4])
		if int(this.MsgSize) == len(buf) {
			this.MsgSender = binary.LittleEndian.Uint32(buf[4:8])
			this.MsgReceiver = binary.LittleEndian.Uint32(buf[8:12])
			this.MsgBody = append(this.MsgBody, buf...)
		} else {
			this.MsgSize = 0
		}
	}
	return this
}

//string -> SimpleMsg
func (this *SimpleMsg) FromString(fromId int, toId int, msg string) *SimpleMsg {
	this.MsgSender = uint32(fromId)
	this.MsgReceiver = uint32(toId)
	dataFrom := make([]byte, 4)
	dataTo := make([]byte, 4)
	dataSize := make([]byte, 4)
	dataBody := []byte(msg)
	this.MsgSize = uint32(len(dataBody) + 12)
	binary.LittleEndian.PutUint32(dataFrom, this.MsgSender)
	binary.LittleEndian.PutUint32(dataTo, this.MsgReceiver)
	binary.LittleEndian.PutUint32(dataSize, this.MsgSize)
	data := []byte{}
	data = append(data, dataSize...)
	data = append(data, dataFrom...)
	data = append(data, dataTo...)
	data = append(data, dataBody...)
	this.MsgBody = data
	return this
}

//SimpleMsg -> []byte
func (this *SimpleMsg) ToData() []byte {
	return this.MsgBody
}

//数据加密 后续实现
func (this *SimpleMsg) EnCode() {
}

//数据解密 后续实现
func (this *SimpleMsg) DeCode() {
}
