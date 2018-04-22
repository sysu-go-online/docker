package util

//********************************************
// Author : huziang
//   包含常用函数
//********************************************

// DealPanic 处理错误
func DealPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func isFileExit() bool {
	return true
}
