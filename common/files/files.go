package files


import "os"

//FileExists 判断文件是否存在
//@path 文件路径
//@return 存在返回true否则返回false
func FileExists(path string) bool {
	_, err :=os.Stat(path)
	if err == nil {
		return true
	}

	return false
}
