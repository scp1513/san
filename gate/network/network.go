// 封装网络相关操作
package network

import (
	"log"
	"time"

	"github.com/scp1513/ec/serial"
)

const (
	statusInvalid = iota
	statusConnecting
	statusRunning
	statusStopping
	statusReconnect
)

var (
	status          = statusInvalid
	sendClientBegin = make(chan struct{}, 0)
	sendShutdown    = make(chan uint32, 1)

	recvSerial             *serial.S
	gameSendSerial         *serial.S
	clientSendSerial       *serial.S
	recvSerialFinish       = make(chan struct{}, 0)
	gameSendSerialFinish   = make(chan struct{}, 0)
	clientSendSerialFinish = make(chan struct{}, 0)
)

// Run
func Run() error {
	startSerial()
	if err := connectGame(); err != nil {
		return err
	}

	status = statusConnecting
	sendClientBegin <- struct{}{}
	startFeedback()
	status = statusRunning
	flag := <-sendShutdown
	status = statusStopping
	shutdown(flag)
	return nil
}

func shutdown(flag uint32) {
	stopFeedback()
	stopListenClient()
	switch flag {
	default:
		log.Printf("Invalid shutdown flag %d. Shutdown normally.", flag)
		fallthrough
	case 1: // 正常关闭
		notifyClientShutdown()
		closeAllClient()
		time.Sleep(time.Second * 5)
		stopSerial()
		waitSerial()
	case 2: // 等待所有客户端连接断开
		for {
			n := GetClientNum()
			if n == 0 {
				break
			}
			log.Println("Wait all client close. Left client num ", n)
			time.Sleep(time.Second * 30)
		}
		time.Sleep(time.Second * 3)
		stopSerial()
		waitSerial()
	case 3: // 强制关闭
		notifyClientShutdown()
		closeAllClient()
		time.Sleep(time.Second * 5)
		stopSerial()
	}
}

func startSerial() {
	recvSerial = serial.New(4096)
	gameSendSerial = serial.New(4096)
	clientSendSerial = serial.New(4096)
	go func() {
		recvSerial.Run()
		<-recvSerialFinish
	}()
	go func() {
		gameSendSerial.Run()
		<-gameSendSerialFinish
	}()
	go func() {
		clientSendSerial.Run()
		<-clientSendSerialFinish
	}()
}

func stopSerial() {
	recvSerial.Stop()
	gameSendSerial.Stop()
	clientSendSerial.Stop()
}

func waitSerial() {
	recvSerialFinish <- struct{}{}
	gameSendSerialFinish <- struct{}{}
	clientSendSerialFinish <- struct{}{}
}
