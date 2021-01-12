package md5

import (
	"crypto/md5"
	"fmt"
	"io"
)

//SignFile 计算文件md5
func SignFile(f io.Reader) string {
	hashMd5 := md5.New()
	io.Copy(hashMd5, f)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}
