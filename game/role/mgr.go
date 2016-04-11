package role

import (
	"github.com/scp1513/san/def"
)

var roles = make(map[uint64]*R)

func Cache(data *def.RoleData) *R {
	if r, ok := roles[data.Rid]; ok {
		r.fromData(data)
		return r
	}
	r := new(R)
	r.fromData(data)
	roles[r.Data.Rid] = r
	return r
}

func Get(rid uint64) *R {
	return roles[rid]
}
