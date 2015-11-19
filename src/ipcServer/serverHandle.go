/*******************************************************************
 *  Copyright(c) 2000-2015 rjcb99
 *  All rights reserved.
 *
 *  文件名称: serverHandle.go
 *  简要描述: 服务控制器
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
	"net"
	"sync"
)

type ServerHandle struct {
	Id   int
	Ip   string
	Port int
	PM   *PostMan
	lock *sync.Mutex
}

//创建服务
func MakeNewServerHandle(id int, ip string, port int, pm *PostMan) *ServerHandle {
	if id < 1 || ip == "" || port < 1 || pm == nil {
		return nil
	}
	return &ServerHandle{
		Id:   id,
		Ip:   ip,
		Port: port,
		PM:   pm,
		lock: new(sync.Mutex),
	}
}

//开始服务
func (this *ServerHandle) Server() {
	connStr := fmt.Sprintf("%s:%d", this.Ip, this.Port)
	l, err := net.Listen("tcp", connStr)
	if err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("[serverHandle:%d] Listen On %s Error: %s", this.Id, connStr, err.Error()))
		return
	}
	SysLog.PutLineAsLog(fmt.Sprintf("[serverHandle:%d] Start On[%s] Successful!", this.Id, connStr))
	for {
		conn, err := l.Accept()
		if err != nil {
			SysLog.PutLineAsLog(fmt.Sprintf("[serverHandle:%d]  Accept Error: %s", this.Id, err.Error()))
			break
		}
		go this.PM.AddConn(conn)
	}
}
