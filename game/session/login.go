package session

import (
	"fmt"
)

func OnLoginSuccess(sid uint64, aid uint32) (*S, error) {
	_, ok := sessions[sid]
	if ok {
		return nil, fmt.Errorf("重复会话ID %d", sid)
	}
	s := createSession(sid, aid)
	sessions[sid] = s
	return s, nil
}
