// 封装网络相关操作
package network

import (
	"github.com/scp1513/ec/serial"
)

var (
	sendShutdown = make(chan uint32, 1)

	recvSerial       *serial.S
	sendSerial       *serial.S
	recvSerialFinish = make(chan struct{}, 0)
	sendSerialFinish = make(chan struct{}, 0)
)

// Run
func Run() error {
	startSerial()
	if err := startListenInner(); err != nil {
		return err
	}

	flag := <-sendShutdown
	shutdown(flag)
	return nil
}

func shutdown(flag uint32) {

}

func startSerial() {
	recvSerial = serial.New(4096)
	sendSerial = serial.New(4096)
	go func() {
		recvSerial.Run()
		<-recvSerialFinish
	}()
	go func() {
		sendSerial.Run()
		<-sendSerialFinish
	}()
}

func stopSerial() {
	recvSerial.Stop()
	sendSerial.Stop()
}

func waitSerial() {
	recvSerialFinish <- struct{}{}
	sendSerialFinish <- struct{}{}
}
