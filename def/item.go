package def

type Item struct {
	ID   uint64 `bson:"id"`   // ID
	Type uint32 `bson:"type"` // 物品类型
}
