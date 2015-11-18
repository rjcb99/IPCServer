package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

//逐行读取文件
func ReadEachLine(path string, dealFunc func([]byte)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadBytes('\n')
		dealFunc(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

//逐行读取文件带终止
func ReadEachLineBreak(path string, dealFunc func([]byte) bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadBytes('\n')
		read_flag := dealFunc(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if read_flag == true {
			return nil
		}
	}
	return nil
}

//一次读取文件
func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

//覆盖保存文件
func SaveFile(path, buf string) bool {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666) //打开文件，不存在就创建，添加模式打开，清空文档，读写打开
	if err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("SaveFile Error : %s", err.Error()))
		return false
	}
	defer f.Close()
	_, err1 := f.WriteString(buf)
	if err1 != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("SaveFile Error : %s", err1.Error()))
		return false
	}
	f.Sync()
	return true
}
