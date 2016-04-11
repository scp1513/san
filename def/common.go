package def

import (
	"fmt"
)

type Pos struct {
	X int `bson:"x"`
	Y int `bson:"y"`
}

func (p *Pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}
