package def

// 属性类型
const (
	PROP_STRENGTH = iota // 力量
	PROP_COUNT
)

// 属性
type Props [PROP_COUNT]int32
