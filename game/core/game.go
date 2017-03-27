package core

import (
	"github.com/lkj01010/goutils/log"
)

const (
	MAX_PLAYER = 5
)

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
		fnCh:      make(chan HandleFunc, GAME_FN_CH_SIZE),
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

func (g *game) H_sessionEnterAck(ss *Session) {
	//log.Debugf("[session <- H_streamSend][userId=%+v]", ss.GetUserId())
	ss.fnCh <- func(ss *Session) {
		ss.H_streamSend(enc_g_enterGameAck(ss.UserId, g.id))
	}
}

func (g *game) H_sessionCloseAck(ss *Session) {
	log.Debugf("[session <- Die][userId=%+v]", ss.GetUserId())
	ss.fnCh <- func(ss *Session) {
		ss.Die()
	}
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

func (g *game) kickSession(ss *Session) {
	//log.Debugf("[kickSession][userId=%+v]", ss.GetUserId())
	if _, ok := g.playerMap[ss.GetUserId()]; ok {
		g.H_sessionLeave(ss)
		return
	} else {
		log.Debugf("[session <- H_close][userId=%+v]", ss.GetUserId())
		ss.fnCh <- func(ss *Session) {
			ss.H_close()
		}
	}
}

func (g *game) H_sessionLeave(ss *Session) {
	//log.Debugf("[game H_sessionLeave()][userId=%+v]", ss.GetUserId())
	for _, p := range g.playerMap {
		p.ss.fnCh <- func(ss *Session) {
			ss.H_streamSend(enc_g_leaveGameNtf(ss.GetUserId()))
		}
	}
	delete(g.playerMap, ss.GetUserId())
	//log.Debugf("[game H_sessionLeave() send H_close()][userId=%+v]", ss.GetUserId())
	log.Debugf("[session <- H_close][userId=%+v]", ss.GetUserId())
	ss.fnCh <- func(ss *Session) {
		ss.H_close()
	}
}

func (g *game) H_enterGame(ss *Session, roomType int32, figure int32) {
	if len(g.playerMap) == MAX_PLAYER {
		panic("gamemgr distribute a full game to player to enter")
	} else {
		if _, ok := g.playerMap[ss.GetUserId()]; ok {
			log.Errorf("[gf_enterGame user exist][id=%+v]", ss.GetUserId())
			g.kickSession(ss)
			return
		}

		player := NewPlayer(ss)

		log.Debugf("[session <- H_enterGame][userId=%+v]", ss.GetUserId())
		ss.fnCh <- func(ss *Session) {
			ss.H_enterGame(g.fnCh)
		}
		g.playerMap[ss.GetUserId()] = player
		for _, p := range g.playerMap {
			if p.ss != ss {
				p.ss.fnCh <- func(ss *Session) {
					ss.H_streamSend(enc_g_enterGameNtf(ss.GetUserId()))
				}
			}
		}
	}
}

func (g *game) h_leaveGame(ss *Session) {
	//log.Debugf("[game H_leaveGame][userId=%+v]", ss.GetUserId())
	//atomic.StoreInt32(&ss.isInGame, 0)
	//close(ss.ToGameCh)

	g.H_sessionLeave(ss)
}