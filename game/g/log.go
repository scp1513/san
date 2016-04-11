package g

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	L         *log.Logger
	logOutput *os.File
)

func InitLogger(name string) (err error) {
	os.Mkdir("log", 0644)
	fn := "log/" + time.Now().Format("20060102150405") + ".log"
	logOutput, err = os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0x644)
	if err != nil {
		return
	}
	L = log.New(io.MultiWriter(os.Stdout, logOutput), name, log.LstdFlags)
	return
}

func ReleaseLogger() error {
	L.SetOutput(os.Stdout)
	return logOutput.Close()
}
