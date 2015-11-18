package utils

import (
	"bytes"
	"fmt"
	log "github.com/cihub/seelog"
	"strconv"
	"strings"
)

func S2i(s string) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		log.Debugf("s2i(%s) err:%s", s, err.Error())
		return 0
	}
	return int(i)
}

func S2i64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		log.Debugf("s2i(%s) err:%s", s, err.Error())
		return 0
	}
	return i
}

func I2s(i int) string {
	return fmt.Sprintf("%d", i)
}

func I2s64(i int64) string {
	return fmt.Sprintf("%d", i)
}

func Ints2s(i []int, sep string) string {
	str := ""
	for _, j := range i {
		if str != "" {
			str += sep
		}
		str += I2s(j)
	}
	return str
}

func S2ints(str, sep string) []int {
	if str == "" {
		return nil
	}
	//解析数据
	nums := []int{}
	strs := strings.Split(str, sep)
	for _, j := range strs {
		nums = append(nums, S2i(j))
	}
	return nums
}

func Ints(ints32 []int32) []int {
	result := make([]int, len(ints32))
	for i, n := range ints32 {
		result[i] = int(n)
	}
	return result
}

func Ints32(ints []int) []int32 {
	result := make([]int32, len(ints))
	for i, n := range ints {
		result[i] = int32(n)
	}
	return result
}

//string -> []byte
func Str2Byte(str string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteString(str)
	return buf.Bytes()
}

func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

func IntsCount(ints []int, n int) int {
	count := 0
	for _, d := range ints {
		if d == n {
			count++
		}
	}
	return count
}

func IntsPos(ints []int, n int) int {
	for i, d := range ints {
		if d == n {
			return i
		}
	}

	return -1
}

func StringPos(strs []string, dst string) int {
	for i, s := range strs {
		if s == dst {
			return i
		}
	}
	return -1
}
