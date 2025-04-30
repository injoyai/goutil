package g

import (
	"github.com/injoyai/base/types"
	"github.com/injoyai/conv"
)

type (
	Var   = conv.Var
	DMap  = conv.Map
	Any   = interface{}
	List  []interface{}
	Bytes = types.Bytes
	M     = Map
)

//========================================Map========================================

type Map types.Maps[string, any]
