package role

import (
	"github.com/golang/protobuf/proto"

	"github.com/scp1513/san/def"
	"github.com/scp1513/san/proto/rsp"
)

type R struct {
	Data *def.RoleData
}

func (r *R) fromData(data *def.RoleData) {
	r.Data = data
}

func (r *R) ToProto() *rsp.RoleData {
	return &rsp.RoleData{
		RoleID: proto.Uint64(r.Data.Rid),
	}
}
