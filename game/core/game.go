package core

import (
	"github.com/lkj01010/goutils/log"

	. "github.com/lkj01010/act-srv/game/com"
)

const (
	MAX_PLAYER = 5
)

type gameHandleFunc func(*game)

type game struct {
	id int32

	fnCh chan HandleFunc
	die chan struct{}

	playerMap map[int32]*Player
}


func init() {
}

func newGame(id int32) (g *game) {
	g = &game{
		id: id,
		fnCh: make(chan HandleFunc, 5),
		die: make(chan struct{}),
		playerMap: make(map[int32]*Player),
	}
	go g.serve()
	return
}

func (g *game) serve() {
	defer func() {
		close(g.fnCh)
	}()

	for {
		select {
		case fn := <-g.fnCh:
			h := fn.(func(g *game))
			h(g)
		case <- g.die:
			return
		}
	}
}

func (g *game) checkEmpty() {

}

func GH_playerLeave(ss *Session, payload []byte) gameHandleFunc {
	return func(g *game) {
		//GameMgr.fnCh <- GMF_gameLeave(ss)
	}
}

func (g *game) h_enterGame(ss *Session, roomType int32, figure int32) {
	log.Debugf("[enterGame][userId=%+v]", ss.GetUserId())
	if len(g.playerMap) == MAX_PLAYER {
		GameMgr.fnCh <- func(gm *gameManager) {
			gm.h_enterGame(ss, roomType, figure)
		}
	} else {
		if _, ok := g.playerMap[ss.GetUserId()]; ok {
			log.Errorf("[gf_enterGame user exist][id=%+v]", ss.GetUserId())
			close(ss.Die)
			return
		}

		player := NewPlayer(ss)
		ss.ToGameCh = g.fnCh
		g.playerMap[ss.GetUserId()] = player
		for _, p := range g.playerMap {
			p.ss.ToAgentCh <- enc_g_enterGameNtf(ss.GetUserId())
		}
	}
}

func (g *game) h_leaveGame(ss *Session) {
	log.Debugf("[leaveGame][userId=%+v]", ss.GetUserId())
	for _, p := range g.playerMap {
		p.ss.ToAgentCh <- enc_g_leaveGameNtf(ss.GetUserId())
	}
	delete(g.playerMap, ss.GetUserId())
	if len(g.playerMap) == 0 {
		close(g.die)
		GameMgr.fnCh <- func(gm *gameManager) {
			gm.h_closeGame(g)
		}
	} else {
		GameMgr.fnCh <- func(gm *gameManager) {
			gm.h_moveGameToVacantList(g)
		}
	}
}
