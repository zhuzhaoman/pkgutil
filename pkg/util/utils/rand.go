package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	gouuid "github.com/nu7hatch/gouuid"
	"github.com/zhuzhaoman/pkgutil/pkg/log"
	"strings"
)

// 数组转，号分割字符串
func Convert(array interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ",", -1)
}

// SignKey 随机生成32的密钥
func SignKey() string {
	var b = make([]byte, 20)
	rand.Read(b)
	mb := md5.Sum(b)
	return fmt.Sprintf("%x", mb[:])
}
func SerialNumber() string {
	u4, err := gouuid.NewV4()
	if err != nil {
		log.Errorf("error: %s", err)
		return ""
	}
	return fmt.Sprintf("%x", u4[:])
}
