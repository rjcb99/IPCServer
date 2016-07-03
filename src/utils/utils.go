package utils

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

const (
	modName   = "[utils]"
	DEBUGFLAG = false
)

//全局参数
var randnum *rand.Rand
var SysLog *MyLog

func init() {
	//init randnum
	seed, dev := time.Since(time.Date(2013, 12, 3, 4, 5, 6, 789, time.UTC)).Nanoseconds(), devRandom()
	if dev != 0 {
		seed *= int64(dev)
	}
	rand.Seed(seed)
	randnum = rand.New(rand.NewSource(time.Now().UnixNano()))
	//init SysLog
	if !Init_SysLog("") {
		fmt.Printf("%s Init_SysLog ERROR!", modName)
	}
}

func Init_SysLog(appName string) bool {
	if SysLog == nil {
		SysLog = MakeNewMyLog(appName+"logs", appName+"_sys.log", 10000000, 5)
	}
	if SysLog == nil {
		return false
	} else {
		return true
	}
}

//get /dev/random uint64
func devRandom() uint64 {
	if file, err := os.Open("/dev/random"); err == nil {
		defer file.Close()
		data := make([]byte, 8)
		if n, e := file.Read(data); n == 8 && e == nil {
			return binary.LittleEndian.Uint64(data)
		} else {
			SysLog.PutLineAsLog(fmt.Sprintf("%s reading /dev/random err: %d, %+v", modName, n, e))
		}
	} else {
		SysLog.PutLineAsLog(fmt.Sprintf("%s open /dev/random err: %s", modName, err.Error()))
	}
	return 0
}

//Random  [m,n] no safe
func GetRandom(m, n int) int {
	if randnum == nil {
		randnum = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	if m > n {
		return 0
	}
	if m == n {
		return m
	}
	return m + randnum.Intn(n-m+1)
}

func Md5hash(s string) string {
	hash := md5.New()
	io.WriteString(hash, s)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

//http://stackoverflow.com/questions/12264789/shuffle-array-in-go
//乱序数组
func Shuffle(data []int) {
	perm := rand.Perm(len(data))
	for i := range data {
		data[i], data[perm[i]] = data[perm[i]], data[i]
	}
}

func SendEmail(subject, body string, recipients []string) {
	to := strings.Join(recipients, "; ")
	title := fmt.Sprintf("Subject: %s\r\nTo: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n", subject, to)
	auth := smtp.PlainAuth(
		"",
		"rjcb99@163.com", //xxxx@xxxx.com
		"**************", //password
		"smtp.163.com",   //smtp.xxx.com
	)
	err := smtp.SendMail(
		"smtp.163.com:25",
		auth,
		"rjcb99@163.com",
		recipients,
		[]byte(title+body),
	)
	if err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("%s send mail '%s' to '%+v' err: %s", modName, subject, recipients, err.Error()))
	} else {
		SysLog.PutLineAsLog(fmt.Sprintf("%s send mail '%s' to '%+v' ok", modName, subject, recipients))
	}
}

func SendUdp(msg string) (int, error) {
	if !DEBUGFLAG {
		return 0, nil
	}
	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8001")
	if err != nil {
		return 0, err
	}
	udp, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return 0, err
	}
	defer udp.Close()
	return udp.Write(Str2Byte(msg))
}
