package container

//********************************************
// Author : huziang
//   将ch当做锁使用，保证信号不会出问题
//********************************************

// Signal 异步信号量
type Signal struct {
	num int
	ch  chan bool
}

// NewSignal 新创建一个信号量指针
func NewSignal(num int) *Signal {
	signal := new(Signal)
	signal.num = num
	return signal
}

// Signal 发送信号
func (signal *Signal) Signal() {
	signal.num--
	signal.ch <- true
	return
}

// Wait 等待信号
func (signal *Signal) Wait() {
	for signal.num > 0 {
		<-signal.ch
	}
}
