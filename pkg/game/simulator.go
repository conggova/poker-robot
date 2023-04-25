package game

import (
	"github.com/conggova/poker-robot/pkg/action"
	"github.com/conggova/poker-robot/pkg/player"
	"github.com/conggova/poker-robot/pkg/strategy"
)

// 载入Player信息
func (*Game) BuildGame(gbi strategy.GameBuildInfo) strategy.Simulator {
	g := newGame([3]player.Player{
		player.NewStrategyPlayer(0, gbi.Strategys[0], strategy.NoCheat, strategy.ShareInterest, nil),
		player.NewStrategyPlayer(1, gbi.Strategys[1], strategy.NoCheat, strategy.ShareInterest, nil),
		player.NewStrategyPlayer(2, gbi.Strategys[2], strategy.NoCheat, strategy.ShareInterest, nil)},
		gbi.OpenCard, false)
	//设置Player与打牌无关的数据
	for i := 0; i <= 2; i++ {
		//有两方是一伙
		if gbi.Partners != [2]int{} && (i == gbi.Partners[0] || i == gbi.Partners[1]) {
			//calcu cheatFlag
			partnerId := gbi.Partners[0]
			if i == gbi.Partners[0] {
				partnerId = gbi.Partners[1]
			}
			partnerFlag := strategy.WithPreer
			if (i+1)%3 == partnerId {
				partnerFlag = strategy.WithNexter
			}
			g.players[i].SetCheatFlag(partnerFlag)
			g.players[i].SetCheatMethod(gbi.CheatMethod)
			g.players[i].SetPartner(g.players[partnerId])
		} else {
			g.players[i].SetCheatFlag(strategy.NoCheat)
		}
	}
	return g
}

// 只载入牌局, partners 0个或2个元素
func (g *Game) RestoreGame(gri strategy.GameRestoreInfo) {
	//要重置GS
	g.gs = newGameStat(g.gs.keepLog)
	//记录之前的两手牌用来判断炸弹, gs.playLog 不用全部加载
	pid := (gri.LastPlayer + 2) % 3
	cnt := len(gri.PlayLogs[pid])
	if cnt != 0 {
		g.gs.putAction(pid, gri.PlayLogs[pid][cnt-1])
		g.llastAction = gri.PlayLogs[pid][cnt-1]
	}
	pid = gri.LastPlayer
	cnt = len(gri.PlayLogs[pid])
	if cnt != 0 {
		g.gs.putAction(pid, gri.PlayLogs[pid][cnt-1])
		g.lastAction = gri.PlayLogs[pid][cnt-1]
	}
	for i := 0; i < 3; i++ {
		g.playLogs[i] = copySlice(gri.PlayLogs[i])
		g.players[i].SetRemainPokerSet(gri.RemainPokerSets[i])
	}
}

// 指定上下文开始游戏
func (g *Game) RunRestoredGame() {
	//看一看最后的出牌者是不是已经出完了
	if g.players[g.gs.lastPlayer].GetRemainPokerSet() == 0 {
		g.gs.endGame([3]action.PokerSet2{
			g.players[0].GetRemainPokerSet(),
			g.players[1].GetRemainPokerSet(),
			g.players[2].GetRemainPokerSet()})
		g.gs.analyze()
		return
	}
	g.enterLoop((g.gs.lastPlayer + 1) % 3)
}

func copySlice(s []action.Action) []action.Action {
	res := make([]action.Action, len(s))
	copy(res, s)
	return res
}
