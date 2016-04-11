package session

import (
	"github.com/scp1513/san/game/network"
)

type S struct {
	ID     uint64
	Aid    uint32
	RoleID uint64
}

func createSession(sid uint64, aid uint32) *S {
	return &S{
		ID:  sid,
		Aid: aid,
	}
}

func (s *S) Send(m network.Message) {
	network.SendClient(s.ID, m)
}
