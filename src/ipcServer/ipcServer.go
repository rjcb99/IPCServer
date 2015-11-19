/*******************************************************************
 *  Copyright(c) 2000-2015 rjcb99
 *  All rights reserved.
 *
 *  文件名称: ipcServer.go
 *  简要描述: 一个简单的信息转发服务
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
	"flag"
	"net"
	"utils"
)

var SysLog *utils.MyLog

func main() {
	ip := flag.String("ip", GetLocalIp(), "ip addr")
	port := flag.Int("port", 8384, "port")
	flag.Parse()
	if !Init_SysLog() {
		println("Init_SysLog() False!")
	}
	if pm := MakeNewPostMan(); pm != nil {
		MakeNewServerHandle(1, *ip, *port, pm).Server()
	}
}

//初始化进程日志
func Init_SysLog() bool {
	if SysLog == nil {
		SysLog = utils.MakeNewMyLog("ipcServerLog", "ipcServer.log", 10000000, 5)
	}
	if SysLog == nil {
		return false
	} else {
		return true
	}
}

//获取本地ip
func GetLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String() //返回第一个
				}
			}
		}
	}
	return "127.0.0.1"
}
