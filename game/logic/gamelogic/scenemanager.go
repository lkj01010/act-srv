package gamelogic

import (
	"container/list"

	"github.com/lkj01010/goutils/log"
)

type sceneManager struct {
	in chan SceneManagerHandler
	// 全部scene
	sceneMap map[int32]*Scene
	// 下一个scene编号
	sceneAcc int32
	// 没有满的房间列表
	sceneListNotFull *list.List
}

var sceneMgr sceneManager

func init() {
	sceneMgr = sceneManager{
		sceneMap:         make(map[int32]*Scene),
		sceneAcc:         0,
		sceneListNotFull: list.New(),
	}
}

func (sm *sceneManager) serve() {
	sm.in = make(chan SceneManagerHandler)
	defer func() {
		close(sm.in)
	}()
	for {
		select {
		case handler := <-sm.in:
			log.Infof("%+v", handler)
			handler(sm)
		}
	}
}

func (sm *sceneManager) Perform(h SceneManagerHandler) {
	sm.in <- h
}

//////////////////////////////////////////////////
// SceneManagerHandler
type SceneManagerHandler func(*sceneManager)

func SMH_playerEnter(player *Player) SceneManagerHandler {
	return func(sm *sceneManager) {
		var scene *Scene
		if sm.sceneListNotFull.Len() == 0 {
			id := (int32)(len(sm.sceneMap))
			scene = newScene(id)

			sm.sceneMap[id] = scene
			sm.sceneListNotFull.PushBack(scene)
		}

		sceneElem := sm.sceneListNotFull.Front()
		scene = sceneElem.Value.(*Scene)

		scene.Perform(SH_playerEnter(player))
	}
}

func SMH_sceneFull(player *Player) SceneManagerHandler {
	return func(sm *sceneManager) {
		log.Warningf("sceneFull try another time, current scene len=%v+, current not full len=%v+",
			len(sm.sceneMap), sm.sceneListNotFull.Len())
		SMH_playerEnter(player)(sm)
	}
}

func SMH_playerLeave(player *Player) SceneManagerHandler {
	return func(sm *sceneManager) {

	}
}

func SMH_killPlayerFromScene(player *Player) SceneManagerHandler {
	return func(sm *sceneManager) {

	}
}

//
//func (sm *sceneManager) newScene() (sceneId int32, sceneIn chan SceneHandler) {
//	id := (int32)(len(sm.sceneMap))
//	scene := newScene(id)
//	sm.sceneMap[id] = scene
//}

func (sm *sceneManager) PlayerEnter(roomType int32, figure int32, userId int32) int32 {
	return 0
}
