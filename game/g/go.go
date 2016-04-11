package g

import (
	"github.com/scp1513/ec/parallel"
)

func init() {
	parallel.SetRecover(recoverFn)
}

func SetGORecv(v bool) {
	parallel.SetGORecv(v)
}

func GO(f func()) {
	parallel.GO(f)
}

func WaitGO() <-chan struct{} {
	return parallel.WaitGO()
}
