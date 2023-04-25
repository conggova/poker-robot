package player

import (
	"github.com/conggova/poker-robot/pkg/action"
	"github.com/conggova/poker-robot/pkg/strategy"
)

type StrategyPlayer struct {
	BasePlayer
	//过牌列表
	strategy strategy.Strategy
}

func NewStrategyPlayer(playerId int, strategy strategy.Strategy, cheatFlag strategy.TCheatFlag, cheatMethod strategy.TCheatMethod, partner Player) *StrategyPlayer {
	return &StrategyPlayer{BasePlayer: BasePlayer{Id: playerId, cheatFlag: cheatFlag, cheatMethod: cheatMethod, partner: partner}, strategy: strategy}
}

func (p *StrategyPlayer) Play(pc strategy.PlayContext) (actionTaken action.Action) {
	defer func() { p.remainPokerSet = p.remainPokerSet.Subtract(actionTaken.PokerSet2()) }()
	//可以通信看到同伙的牌，始终是明牌
	if (p.cheatFlag == strategy.WithNexter || p.cheatFlag == strategy.WithPreer) && p.cheatMethod == strategy.CommInSecret {
		if !pc.GameOpenCard { //靠partner来取得上下家牌信息
			if p.cheatFlag == strategy.WithPreer {
				pc.PreerRemainPokerSet = p.partner.GetRemainPokerSet()
				pc.NexterRemainPokerSet = pc.OthersRemainPokerSet.Subtract(pc.PreerRemainPokerSet)
			} else {
				pc.NexterRemainPokerSet = p.partner.GetRemainPokerSet()
				pc.PreerRemainPokerSet = pc.OthersRemainPokerSet.Subtract(pc.NexterRemainPokerSet)
			}
		}
		return p.strategy.MakeDecision(strategy.StrategyContext{
			RemainPokerSet: p.remainPokerSet,
			CheatFlag:      p.cheatFlag,
			CheatMethod:    p.cheatMethod,
			OpenCard4Me:    true,
			PlayContext:    pc,
		})
	} else {
		return p.strategy.MakeDecision(strategy.StrategyContext{
			RemainPokerSet: p.remainPokerSet,
			CheatFlag:      p.cheatFlag,
			CheatMethod:    p.cheatMethod,
			OpenCard4Me:    pc.GameOpenCard,
			PlayContext:    pc,
		})
	}
}
