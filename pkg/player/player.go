package player

import (
	"github.com/conggova/poker-robot/pkg/action"
	"github.com/conggova/poker-robot/pkg/strategy"
)

type Player interface {
	Play(strategy.PlayContext) action.Action
	GetId() int
	SetRemainPokerSet(action.PokerSet2)
	GetRemainPokerSet() action.PokerSet2
	//与作弊相关的
	SetCheatFlag(strategy.TCheatFlag)
	SetCheatMethod(strategy.TCheatMethod)
	SetPartner(Player)
}

type BasePlayer struct {
	Id             int
	cheatFlag      strategy.TCheatFlag
	cheatMethod    strategy.TCheatMethod
	partner        Player
	remainPokerSet action.PokerSet2
}

func (p BasePlayer) GetId() int {
	return p.Id
}

func (p *BasePlayer) SetRemainPokerSet(ps action.PokerSet2) {
	p.remainPokerSet = ps
}

func (p BasePlayer) GetRemainPokerSet() action.PokerSet2 {
	return p.remainPokerSet
}

func (p *BasePlayer) SetCheatFlag(f strategy.TCheatFlag) {
	p.cheatFlag = f
}

func (p *BasePlayer) SetCheatMethod(t strategy.TCheatMethod) {
	p.cheatMethod = t
}

func (p *BasePlayer) SetPartner(t Player) {
	p.partner = t
}

func copySlice(s []action.Action) []action.Action {
	res := make([]action.Action, len(s))
	copy(res, s)
	return res
}
