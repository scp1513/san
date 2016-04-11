package session

func Logout(sid uint64) {
	delete(sessions, sid)
}
