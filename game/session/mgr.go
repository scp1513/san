package session

var sessions = make(map[uint64]*S)

func Get(sid uint64) *S {
	return sessions[sid]
}

func Count() int {
	return len(sessions)
}
