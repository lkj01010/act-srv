package gamelogic

import "github.com/lkj01010/goutils/log"

const (
	MAX_PLAYER = 5
)

type Scene struct {
	id              int32
	in              chan SceneHandler
	playerMap       map[int32]*Player
	playerAcc       int
	itemList        []item
	itemActiveCount int
}

func newScene(id int32) (scene *Scene) {
	scene = &Scene{
		id:              id,
		playerMap:       make(map[int32]*Player),
		playerAcc:       0,
		itemList:        make([]item, 0),
		itemActiveCount: 0,
	}
	go scene.serve()
	return
}

func (s *Scene) serve() {
	s.in = make(chan SceneHandler, 5)
	defer func() {
		close(s.in)
	}()

	for {
		select {
		case handler := <-s.in:
			log.Infof("%+v", handler)
			handler(s)
		}
	}
}

func (s *Scene) Perform(h SceneHandler) {
	s.in <- h
}

type SceneHandler func(*Scene)

func SH_playerEnter(player *Player) SceneHandler {
	return func(s *Scene) {
		if len(s.playerMap) == MAX_PLAYER {
			sceneMgr.Perform(SMH_sceneFull(player))
		} else {
			s.playerAcc++
			s.playerMap[s.playerAcc] = player
		}

		for _, v := range s.playerMap {
			v.Send_player_enter_game(player)
		}
	}
}

func SH_playerLeave(player *Player) SceneHandler {
	return func(s *Scene) {

		sceneMgr.Perform(SMH_playerLeave(player))
	}
}

//////////////////////////////////////////////////
type itemType int

const (
	gem itemType = iota
	heart
)

type item struct {
	ty itemType
}
