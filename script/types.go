package script

import _ "embed"

var (
	//go:embed redis/interaction_incr_cnt.lua
	LuaInteractionIncrCnt string
)
