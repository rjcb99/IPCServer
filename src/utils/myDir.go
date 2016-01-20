package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/*获取当前文件执行的路径*/
func GetExePath() string {
	file, _ := exec.LookPath(os.Args[0])    //执行文件的路径
	path, _ := filepath.Abs(file)           //补全绝对路径
	splitstring := strings.Split(path, "/") //切割掉执行文件名
	size := len(splitstring)
	ret := "/"
	if size >= 1 {
		splitstring[size-1] = ""
		ret = strings.Join(splitstring, "/")
	}
	return ret
}

/*获取当前文件的路径*/
func GetExeFilePath() string {
	file, _ := exec.LookPath(os.Args[0]) //执行文件的路径
	path, _ := filepath.Abs(file)        //补全绝对路径
	return path
}

/*创建目录*/
func MkDir(path string) bool {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		//SysLog.PutLineAsLog(fmt.Sprintf("MkDir(%s) Error:%s", path, err.Error()))
		return false
	}
	return true
}

/*删除目录*/
func RmDir(path string) bool {
	err := os.RemoveAll(path)
	if err != nil {
		//SysLog.PutLineAsLog(fmt.Sprintf("RmDir(%s) Error:%s", path, err.Error()))
		return false
	}
	return true
}

/*文件或文件夹是否存在*/
func IsExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if !os.IsExist(err) {
			return false
		}
	}
	return true
}

/*文件是否存在*/
func IsFile(file string) bool {
	f, err := os.Stat(file)
	if err != nil {
		//SysLog.PutLineAsLog(fmt.Sprintf("IsFile(%s) Error:%s", file, err.Error()))
		return false
	}
	return !f.IsDir()
}

/*获得文件大小*/
func FileSize(path string) int64 {
	f, err := os.Stat(path)
	if err != nil {
		//SysLog.PutLineAsLog(fmt.Sprintf("FileSize(%s) Error:%s", path, err.Error()))
		return 0
	}
	return f.Size()
}

/*得到文件的修改时间*/
func FileModifyTime(path string) int64 {
	f, err := os.Stat(path)
	if err != nil {
		//SysLog.PutLineAsLog(fmt.Sprintf("FileModifyTime(%s) Error:%s", path, err.Error()))
		return 0
	}
	return f.ModTime().Unix()
}

/*删除文件*/
func RmFile(path string) error {
	return os.Remove(path)
}

/*重命名文件*/
func ReNameFile(old_path string, new_path string) error {
	return os.Rename(old_path, new_path)
}

/*统一格式的时间string*/
func GetTime() string {
	return time.Now().String()[:19]
}
