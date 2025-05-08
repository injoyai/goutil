package g

import (
	"github.com/injoyai/base/types"
	"github.com/injoyai/conv"
)

type (
	Var   = conv.Var
	DMap  = conv.Map
	Any   = interface{}
	Bytes = types.Bytes
	M     = Map
	List  types.List[any]
	Map   types.Map[string, any]
	Maps  types.Maps[string, any]
)
