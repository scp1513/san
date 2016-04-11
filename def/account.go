package def

// AccountData 账号数据
type AccountData struct {
	Aid     uint32   `bson:"_id"`
	RoleIDs []uint64 `bson:"role_ids"`
}

// AccountDebugData 账号调试数据
type AccountDebugData struct {
	Account  string `bson:"account"`
	Password string `bson:"password"`
}
