package core

import (
	"github.com/lkj01010/goutils/log"

	. "github.com/lkj01010/act-srv/game/com"
	"sort"
)

type vacantGameList struct {
	gl []*game
}

// 插入和删除有序的情况下，貌似用不着sort()
func (v vacantGameList) Len() int           { return len(v.gl) }
func (v vacantGameList) Less(i, j int) bool { return len(v.gl[i].playerMap) < len(v.gl[j].playerMap) }
func (v vacantGameList) Swap(i, j int)      { v.gl[i], v.gl[j] = v.gl[j], v.gl[i] }
func (v vacantGameList) Search(vCount int) int {
	return sort.Search(len(v.gl), func(i int) bool {
		return len(v.gl[i].playerMap) == vCount
	})
}
func (v *vacantGameList) Add(g *game) {
	vCount := MAX_PLAYER - len(g.playerMap)
	index := v.Search(vCount)

	gl1 := append(v.gl[:index], g)
	gl2 := append(gl1, v.gl[index:]...)
	v.gl = gl2
}

func (v *vacantGameList) Remove(g *game) {
	vCount := MAX_PLAYER - len(g.playerMap)
	index := v.Search(vCount)
	v.gl = append(v.gl[:index], v.gl[index:]...)
}

func (v *vacantGameList) Update(g *game) {
	v.Remove(g)
	v.Add(g)
}

func (v *vacantGameList) End() *game {
	g := v.gl[len(v.gl)-1]
	return g
}

//////////////////////////////////////////////////
type gameManager struct {
	msgSsCh chan *MsgWithSession
	fnCh    chan HandleFunc
	// 全部scene
	gameMap map[int32]*game
	// 下一个scene编号
	gameAcc int32
	// 没有满的房间列表
	//gameListNotFull *list.List
	vgList *vacantGameList
}

var GameMgr gameManager

func init() {
	GameMgr = gameManager{
		fnCh:    make(chan HandleFunc, 10),
		gameMap: make(map[int32]*game),
		gameAcc: 0,
		//gameListNotFull: list.New(),
		vgList: &vacantGameList{make([]*game, 0, 5)},
	}
	go GameMgr.serve()
}

func (gm *gameManager) serve() {
	defer func() {
		//close(gm.msgSsCh)
		close(gm.fnCh)
	}()
	for {
		select {
		//case handler := <-gm.in:
		//	log.Infof("%+v", handler)
		//	handler(gm)

		//case msgWithSession := <-gm.msgSsCh:
		//	if h, ok := gm.handlers[msgWithSession.Cmd]; !ok {
		//		log.Error("not found cmd=%+v", msgWithSession.Cmd)
		//	} else {
		//		h(msgWithSession.Ss, msgWithSession.Payload)
		//	}
		case fn := <-gm.fnCh:
			h := fn.(func(gm *gameManager))
			h(gm)
		}
	}
}

//////////////////////////////////////////////////
// gameManagerHandler

//func GMH_playerEnter(ss *Session, roomType int32, figure int32, userId int32) gameManagerHandler {
//	return func(gm *gameManager) {
//		var g *game
//		if gm.gameListNotFull.Len() == 0 {
//			id := (int32)(len(gm.gameMap))
//			g = newGame(id)
//
//			gm.gameMap[id] = g
//			gm.gameListNotFull.PushBack(g)
//		}
//
//		sceneElem := gm.gameListNotFull.Front()
//		g = sceneElem.Value.(*game)
//
//		g.msgCh <- GH_sessionEnter(ss, roomType, figure, userId)
//	}
//}
//
//func GMH_gameFull(roomType int32, figure int32, userId int32) gameManagerHandler {
//	return func(gm *gameManager) {
//		log.Warningf("sceneFull try another time, current scene len=%v+, current not full len=%v+",
//			len(gm.gameMap), gm.gameListNotFull.Len())
//		GMH_playerEnter(roomType, figure, userId)(gm)
//	}
//}
//

//
//func SMH_killPlayerFromScene(player *Player) gameManagerHandler {
//	return func(gm *gameManager) {
//
//	}
//}

//
//func (gm *sceneManager) newScene() (sceneId int32, sceneIn chan SceneHandler) {
//	id := (int32)(len(gm.gameMap))
//	scene := newScene(id)
//	gm.gameMap[id] = scene
//}

func (gm *gameManager) PlayerEnter(roomType int32, figure int32, userId int32) int32 {
	return 0
}

//////////////////////////////////////////////////

func (gm *gameManager) h_enterGame(ss *Session, roomType int32, figure int32) {
	var g *game
	if gm.vgList.Len() == 0 {
		id := (int32)(len(gm.gameMap))
		g = newGame(id)
		gm.gameMap[id] = g
		gm.vgList.Add(g)
	}
	g = gm.vgList.End()
	log.Infof("[get a game][id=%+v][player count=%+v][total game count=%+v]", g.id, len(g.playerMap), len(gm.gameMap))

	g.fnCh <- func(g *game) {
		g.h_enterGame(ss, roomType, figure)
	}
}

func (gm *gameManager) h_closeGame(g *game) {
	if len(g.playerMap) > 0 {
		log.Error("close game when it is not empty")
	}
	gm.vgList.Remove(g)
	delete(gm.gameMap, g.id)
}

func (gm *gameManager) h_moveGameToVacantList(g *game) {
	gm.vgList.Update(g)
}
