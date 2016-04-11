package def

// SelectRoleInfo 选择角色界面的角色信息
type SelectRoleInfo struct {
	Rid  uint64 `bson:"_id"`
	Name string `bson:"name"`
}

// RoleData 数据库里的角色数据
type RoleData struct {
	Rid    uint64  `bson:"_id"`    // 角色ID
	Aid    uint32  `bson:"aid"`    // 账号ID
	Name   string  `bson:"name"`   // 昵称
	Coin   uint32  `bson:"coin"`   // 游戏币
	Gold   uint32  `bson:"gold"`   // 金币
	Exp    uint32  `bson:"exp"`    // 经验
	Level  uint16  `bson:"level"`  // 等级
	Props  Props   `bson:"props"`  // 基本属性
	Items  []Item  `bson:"items"`  // 物品
	Equips []Equip `bson:"equips"` // 装备
}
