package game

import . "github.com/lkj01010/act-srv/game/registry"

var Ipc Registry

func init() {
	Ipc.Init()
}