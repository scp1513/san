package cnf

// 常量表

import (
	"encoding/csv"
	"fmt"
	"strconv"
)

// ConstCfg 常量配置
type ConstCfg struct {
	ID    int   `csv:"0,key"` // id
	Value int64 `csv:"1"`     // 常量值
}

// 常量表配置
var constCfgMap map[int]*ConstCfg

func loadConstTable(r *csv.Reader) error {
	m := make(map[int]*ConstCfg)
	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		if len(record) != 3 {
			return fmt.Errorf("const row len error %d", len(record))
		}

		id, err := strconv.ParseInt(record[0], 10, 32)
		if err != nil {
			return err
		}
		val, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return err
		}
		m[int(id)] = &ConstCfg{int(id), val}
	}

	constCfgMap = m
	return nil
}

func GetConstInt(id int) (int, bool) {
	v, ok := constCfgMap[id]
	if !ok {
		return 0, false
	}
	return int(v.Value), true
}

func GetConstUint(id int) (uint, bool) {
	v, ok := constCfgMap[id]
	if !ok {
		return 0, false
	}
	return uint(v.Value), true
}

func GetConstInt8(id int) (int8, bool) {
	v, ok := constCfgMap[id]
	if !ok {
		return 0, false
	}
	return int8(v.Value), true
}

func GetConstUint8(id int) (uint8, bool) {
	v, ok := constCfgMap[id]
	if !ok {
		return 0, false
	}
	return uint8(v.Value), true
}

func GetConstInt16(id int) (int16, bool) {
	v, ok := constCfgMap[id]
	if !ok {
		return 0, false
	}
	return int16(v.Value), true
}

func GetConstUint16(id int) (uint16, bool) {
	v, ok := constCfgMap[id]
	if !ok {
		return 0, false
	}
	return uint16(v.Value), true
}

func GetConstInt32(id int) (int32, bool) {
	v, ok := constCfgMap[id]
	if !ok {
		return 0, false
	}
	return int32(v.Value), true
}

func GetConstUint32(id int) (uint32, bool) {
	v, ok := constCfgMap[id]
	if !ok {
		return 0, false
	}
	return uint32(v.Value), true
}

func GetConstInt64(id int) (int64, bool) {
	v, ok := constCfgMap[id]
	if !ok {
		return 0, false
	}
	return v.Value, true
}

func GetConstUint64(id int) (uint64, bool) {
	v, ok := constCfgMap[id]
	if !ok {
		return 0, false
	}
	return uint64(v.Value), true
}
