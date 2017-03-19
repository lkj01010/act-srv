package core

import (
	. "github.com/lkj01010/act-srv/game/com"
	. "github.com/lkj01010/act-srv/com"
	"github.com/lkj01010/goutils/log"
)

func Dispatch(ss *Session, cmd int16, payload []byte) {
	dest := CmdSendTo(cmd)
	if dest == DestGameMgr {
		GameMgr.fnCh <- gameMgrDecoders[cmd](ss, payload)
	} else if dest == DestGame {
		if ss.ToGameCh == nil {
			//ss.Die <- struct{}{}
			log.Error("dispatch to game when ToGameCh==nil")
		} else {
			ss.ToGameCh <- gameDecoders[cmd](ss, payload)
		}
	} else {
		log.Errorf("[cmd not found][cmd=%+v]", cmd)
	}
}

//func Kick(ss *Session) {
//	// note: 犹豫是异步，此时可能ss已经die了，但为了使流程简单化，不用flag来判断，而是就让ss往前继续走，
//	// 反正走步了几步，要是会出问题再用启用flag
//	//if ss.ToGameCh != nil {
//	//	ss.ToGameCh <- dec_g_leaveGame(ss, nil)
//	//}
//	ss.Flag |= SessKickedOut
//}