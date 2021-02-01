package md5

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

//SignFile 计算文件md5
func SignFile(f io.Reader) string {
	hashMd5 := md5.New()
	io.Copy(hashMd5, f)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

//Sign md5加密
func Sign(origin string) string {
	hashMd5 := md5.New()
	_, _ = io.WriteString(hashMd5, origin)
	return strings.ToLower(fmt.Sprintf("%x", hashMd5.Sum(nil)))
}
