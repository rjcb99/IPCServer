package main

import (
	"flag"
	"net"
	"utils"
)

/*  说明:一个简单的信息转发服务
simpleMsg.go    : 信息格式 信息最大长度65536byte
msgBox.go       : 信箱，存放信息，可以上线后接收，信箱有容量限制，默认2000，牺牲旧数据
postMan.go      : 邮递员，负责开辟信箱和信息的分发
serverHandle.go : 用户上线通知邮递员
*/
var SysLog *utils.MyLog

func main() {
	//初始化参数
	//	ip := flag.String("ip", GetLocalIp(), "ip addr")
	//	port := flag.Int("port", 8384, "port")
	flag.Parse()
	Init_SysLog()
	//pm := MakeNewPostMan()
	//MakeNewServerHandle(1, *ip, *port, pm).Server()
}

//初始化进程日志
func Init_SysLog() bool {
	if SysLog == nil {
		//		SysLog = utils.MakeNewMyLog("IpcServer_logs", "Ipc_sys.log", 10000000, 5)
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
	if err != nil {
		//		SysLog.PutLineAsLog(fmt.Sprintf("GetLocalIp Error:%s So Use 127.0.0.1", err.Error()))
		return "127.0.0.1"
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String() //返回第一个
			}
		}
	}
	return "127.0.0.1"
}
