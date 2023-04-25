package game

import (
	"github.com/conggova/poker-robot/pkg/action"
	"github.com/conggova/poker-robot/pkg/player"
	"github.com/conggova/poker-robot/pkg/strategy"
)

// Game要满足Simulator接口
type Game struct {
	players     [3]player.Player
	playLogs    [3][]action.Action
	lastAction  action.Action
	llastAction action.Action
	gs          GameStat //牌局统计
	openCard    bool     //这局游戏是否明牌
}

func newGame(players [3]player.Player, openCard bool, keepLog bool) *Game {
	return &Game{players: players, gs: newGameStat(keepLog), openCard: openCard}
}

// 返回指定玩家的收益
func (g *Game) GetProfits() [3]int {
	return g.gs.getProfits()
}

func (g *Game) RunGameWithRandomBeginning() {
	//随机发牌
	pokerSets := strategy.DistCards()
	//红3先出
	var headPlayer int
	for i := 0; i < 3; i++ {
		if pokerSets[i]&1 == 1 {
			headPlayer = i
		}
	}
	//转为set2
	pokerSet2s := [3]action.PokerSet2{}
	for i := 0; i < 3; i++ {
		pokerSet2s[i] = pokerSets[i].PokerSet2()
	}
	for i := 0; i < 3; i++ {
		g.players[i].SetRemainPokerSet(pokerSet2s[i])
	}
	g.enterLoop(headPlayer)
}

func (g *Game) enterLoop(currentPlayer int) {
	for {
		var a action.Action
		//如果是明牌，会给Player另外两家的牌形信息
		if g.openCard {
			preerPS := g.players[(currentPlayer+2)%3].GetRemainPokerSet()
			nexterPS := g.players[(currentPlayer+1)%3].GetRemainPokerSet()
			a = g.players[currentPlayer].Play(strategy.PlayContext{
				PreerAction:          g.lastAction,
				NexterAction:         g.llastAction,
				PreerPkCnt:           preerPS.PokerCount(),
				NexterPkCnt:          nexterPS.PokerCount(),
				GameOpenCard:         g.openCard,
				OthersRemainPokerSet: preerPS.CombineWith(nexterPS),
				PreerRemainPokerSet:  preerPS,
				NexterRemainPokerSet: nexterPS,
				PlayLog:              g.playLogs[currentPlayer],
				PreerPlayLog:         g.playLogs[(currentPlayer+2)%3],
				NexterPlayLog:        g.playLogs[(currentPlayer+1)%3],
			})

		} else {
			preerPS := g.players[(currentPlayer+2)%3].GetRemainPokerSet()
			nexterPS := g.players[(currentPlayer+1)%3].GetRemainPokerSet()
			a = g.players[currentPlayer].Play(strategy.PlayContext{
				PreerAction:          g.lastAction,
				NexterAction:         g.llastAction,
				PreerPkCnt:           preerPS.PokerCount(),
				NexterPkCnt:          nexterPS.PokerCount(),
				GameOpenCard:         g.openCard,
				OthersRemainPokerSet: preerPS.CombineWith(nexterPS),
				PreerRemainPokerSet:  0,
				NexterRemainPokerSet: 0,
				PlayLog:              g.playLogs[currentPlayer],
				PreerPlayLog:         g.playLogs[(currentPlayer+2)%3],
				NexterPlayLog:        g.playLogs[(currentPlayer+1)%3],
			})
		}
		g.playLogs[currentPlayer] = append(g.playLogs[currentPlayer], a)
		g.gs.putAction(currentPlayer, a)
		//已经出完
		if g.players[currentPlayer].GetRemainPokerSet() == 0 {
			//结束游戏
			g.gs.endGame([3]action.PokerSet2{
				g.players[0].GetRemainPokerSet(),
				g.players[1].GetRemainPokerSet(),
				g.players[2].GetRemainPokerSet()})
			break
		}
		g.llastAction = g.lastAction
		g.lastAction = a
		//转到下家
		currentPlayer = (currentPlayer + 1) % 3
	}
	//计算收益
	g.gs.analyze()
}
