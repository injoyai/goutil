package g

import (
	"github.com/injoyai/base/types"
	"github.com/injoyai/conv"
)

type (
	Var   = conv.Var
	DMap  = conv.Map
	Any   = any
	Bytes = types.Bytes
	M     = Map
	Map   = types.Map[string, any]
	Maps  = types.Maps[string, any]
	List  = types.List[any]
)
