package gen

import (
	"strconv"
	"time"

	"github.com/scp1513/san/common"
)

func RandAccount() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + common.RandString(8)
}

func RandPassword() string {
	return common.RandString(16)
}

func RandSalt() string {
	return common.RandString(8)
}

func RandToken() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + common.RandString(20)
}
