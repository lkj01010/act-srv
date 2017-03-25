package core

import (
	"github.com/lkj01010/goutils/log"
	"sync/atomic"
)

const (
	MAX_PLAYER = 5
)

type gameHandleFunc func(*game)

type game struct {
	id int32

	fnCh chan HandleFunc
	die  chan struct{}

	playerCount int32 // gamemgr用来做流量控制的标示
	playerMap   map[int32]*Player
}

func init() {
}

func newGame(id int32) (g *game) {
	g = &game{
		id:        id,
		fnCh:      make(chan HandleFunc, 5),
		die:       make(chan struct{}),
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
		case <-g.die:
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
		panic("gamemgr distribute a full game to player to enter")
	} else {
		if _, ok := g.playerMap[ss.GetUserId()]; ok {
			log.Errorf("[gf_enterGame user exist][id=%+v]", ss.GetUserId())
			close(ss.DieCh)
			return
		}

		player := NewPlayer(ss)
		ss.ToGameCh = g.fnCh
		atomic.StoreInt32(&ss.isInGame, 1)
		g.playerMap[ss.GetUserId()] = player
		for _, p := range g.playerMap {
			p.ss.ToAgentCh <- enc_g_enterGameNtf(ss.GetUserId())
		}
	}
}

func (g *game) h_leaveGame(ss *Session) {
	log.Debugf("[leaveGame][userId=%+v]", ss.GetUserId())

	atomic.StoreInt32(&ss.isInGame, 0)
	close(ss.ToGameCh)

	for _, p := range g.playerMap {
		if p.ss == ss {
			if atomic.LoadInt32(&p.ss.isDie) == 1 {
				close(p.ss.DieCh)
			} else {
				p.ss.ToAgentCh <- enc_g_leaveGameNtf(ss.GetUserId())
			}
		} else {
			p.ss.ToAgentCh <- enc_g_leaveGameNtf(ss.GetUserId())
		}
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
