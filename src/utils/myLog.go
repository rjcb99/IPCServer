package utils

import (
	"fmt"
	"os"
	"sync"
)

const (
	FileClose = iota
	FileOpen
)

type MyLog struct {
	state    int
	lock     *sync.RWMutex
	FileDir  string
	FileName string
	Path     string
	F        *os.File
	MaxSize  int64
	MaxFiles int
}

//在当前执行目录下创建 fdir 目录，然后创建 fdir 中的文件fname 需要手工释放
func MakeNewMyLog(fdir string, fname string, max_size int64, max_file int) *MyLog {
	if !IsExist(GetExePath() + fdir) { //文件夹不存在，就地创建
		MkDir(GetExePath() + fdir)
	}
	mf := &MyLog{
		state:    FileClose,
		lock:     new(sync.RWMutex),
		FileDir:  fdir,
		FileName: fname,
		Path:     GetExePath() + fdir + "/" + fname,
		F:        nil,
		MaxSize:  max_size,
		MaxFiles: max_file,
	}
	var err error
	mf.F, err = os.OpenFile(mf.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666) //打开文件，不存在就创建，添加模式打开，读写打开
	if err != nil {
		fmt.Sprintf("MakeNewMyLog Error:%s", err.Error())
		return nil
	}
	mf.state = FileOpen
	return mf
}

func (this *MyLog) Release() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.state == FileOpen {
		this.F.Close()
		this.F = nil
		this.state = FileClose
	}
}

//插入一行 非安全
func (this *MyLog) putLine(buf string) bool {
	if this.state == FileClose {
		return false
	}
	buf = buf + "\n"
	_, err := this.F.WriteString(buf)
	if err != nil {
		fmt.Sprintf("putLine(%s) Error:%s", buf, err.Error())
		return false
	}
	return this.checkBak()
}

//插入一行 安全
func (this *MyLog) PutLine(buf string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.putLine(buf)
}

//插入一行日志 安全
func (this *MyLog) PutLineAsLog(buf string) bool {
	if this != nil {
		return this.PutLine("[" + GetTime() + "]" + " " + buf)
	}
	return false
}

//插入16进制日志 安全
func (this *MyLog) PutHexAsLog(msg []byte, len int) bool {
	ASCII := make([]byte, 18)
	VALUE := make([]byte, 34)
	HEAD_ADR := make([]byte, 5)
	if len > 20000000 {
		return false
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	ASCII[8] = '-'
	ASCII[17] = 0x00
	VALUE[16] = '-'
	VALUE[33] = 0x00
	this.putLine(GetTime() + " ")
	for i := 0; i < len; i++ {
		if i%16 < 8 {
			str := fmt.Sprintf("%02x", msg[i])
			VALUE[(i%16)*2] = str[0]
			VALUE[(i%16)*2+1] = str[1]
			if msg[i] < 32 || msg[i] > 126 {
				ASCII[i%16] = '.'
			} else {
				ASCII[i%16] = msg[i]
			}
		} else {
			str := fmt.Sprintf("%02x", msg[i])
			VALUE[(i%16)*2+1] = str[0]
			VALUE[(i%16)*2+2] = str[1]
			if msg[i] < 32 || msg[i] > 126 {
				ASCII[i%16] = '.'
			} else {
				ASCII[i%16] = msg[i]
			}
		}
		if i == len-1 {
			HEAD_ADR[4] = 0x00
			str := fmt.Sprintf("%003X0", i/16)
			str = str + "[" + string(VALUE) + "] [" + string(ASCII) + "]"
			this.putLine(str)
		} else if i > 0 && i%16 == 15 {
			HEAD_ADR[4] = 0x00
			str := fmt.Sprintf("%003X0", i/16)
			str = str + "[" + string(VALUE) + "] [" + string(ASCII) + "]"
			this.putLine(str)
			ASCII = make([]byte, 18)
			VALUE = make([]byte, 34)
			ASCII[8] = '-'
			ASCII[17] = 0x00
			VALUE[16] = '-'
			VALUE[33] = 0x00
		}
	}
	return true
}

//根据参数检查当前文件是否已经超过了规定的大小
func (this *MyLog) checkBak() bool {
	if this.state == FileClose {
		return false
	}
	//检查文件是否已经超长
	if FileSize(this.Path) < this.MaxSize {
		return true
	}
	/*超长处理*/
	//关闭当前文件
	this.F.Close()
	this.F = nil
	var err error
	if this.MaxFiles <= 1 { //只允许一个文件
		this.F, err = os.OpenFile(this.Path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666) //打开文件，不存在就创建，添加模式打开，清空文档，读写打开
		if err != nil {
			fmt.Sprintf("CheckBak Error:%s", err.Error())
			return false
		}
		return true
	}
	//多个文件就切换最老的文件
	flag := false
	i := 1
	for ; i < this.MaxFiles; i++ { //查找最旧日志名
		if IsFile(this.Path + ".bak_" + I2s(i)) {
			continue
		} else {
			flag = true
			break
		}
	}
	if flag == true {
		ReNameFile(this.Path, this.Path+".bak_"+I2s(i))
	} else {
		RmFile(this.Path + ".bak_1")        //删除备份1
		for i = 2; i < this.MaxFiles; i++ { //剩余备份文件重命名
			ReNameFile(this.Path+".bak_"+I2s(i), this.Path+".bak_"+I2s(i-1))
		}
		ReNameFile(this.Path, this.Path+".bak_"+I2s(i-1)) //当前文件重命名
	}
	this.F, err = os.OpenFile(this.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666) //打开文件，不存在就创建，添加模式打开，读写打开
	if err != nil {
		fmt.Sprintf("CheckBak Error:%s", err.Error())
		return false
	}
	return true
}
