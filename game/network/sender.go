package network

func Send(connID uint, m Message) {
	sendSerial.Post(func() { sendHandler(connID, m) })
}

// SendClient 发送给客户端
func SendClient(sid uint64, m Message) {
	sendSerial.Post(func() { proxyHandler(sid, m) })
}
