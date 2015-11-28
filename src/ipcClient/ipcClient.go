/*******************************************************************
 *  Copyright(c) 2000-2015 rjcb99
 *  All rights reserved.
 *
 *  文件名称: ipcClient.go
 *  简要描述: 一个简单的通过ipcServer通讯的DEMO
 *
 *  创建日期: 2015-11-27
 *  作者: ChenBo
 *  说明:
 *
 *  修改日期: 2015-11-28
 *  作者: ChenBo
 *  说明:
 ******************************************************************/
package main

import (
	"flag"
	"fmt"
	"net"
	"time"
	"utils"
)

var SysLog *utils.MyLog
var Mq *utils.MsgBox

func main() {
	id := flag.Int("id", 1, "client_id")
	ip := flag.String("ip", GetLocalIp(), "ip addr")
	port := flag.Int("port", 8384, "port")
	flag.Parse()
	if !Init_SysLog() {
		println("Init_SysLog() False!")
		return
	}
	//初始化全局Mq
	if Init_Mq(*id, *ip, *port) == false {
		SysLog.PutLineAsLog(fmt.Sprintf("error Init_Mq(%d,%s,%d) : in NN ", *id, *ip, *port))
		return
	}
	for i := 1; i < 3600000; i++ {
		Mq.SendMsg(*id, "北国风光，千里冰封，万里雪飘。望长城内外，惟余莽莽；大河上下，顿失滔滔。山舞银蛇，原驰蜡象，欲与天公试比高。须晴日，看红装素裹，分外妖娆。江山如此多娇，引无数英雄竞折腰。惜秦皇汉武，略输文采；唐宗宋祖，稍逊风骚。一代天骄，成吉思汗，只识弯弓射大雕。俱往矣，数风流人物，还看今朝。")
		time.Sleep(time.Millisecond * 1)
	}
	SysLog.PutLineAsLog("ipcClient Exit!")
}

//初始化进程日志
func Init_SysLog() bool {
	if SysLog == nil {
		SysLog = utils.MakeNewMyLog("ipcClientLog", "ipcClient.log", 10000000, 5)
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

//初始化mq
func Init_Mq(id int, ip string, port int) bool {
	Mq = utils.MakeNewMsgBox(id, ip, port)
	return Mq.MakeConn()
}
