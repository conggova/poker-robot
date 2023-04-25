package strategy

import (
	"github.com/conggova/poker-robot/pkg/action"
)

// 仿真策略
type SimulateStrategy struct {
	distCardTimes        int //进行发牌次数
	simulateTimesPerDist int //每次发牌进行的simulate次数
	simulator            Simulator
	strategyWhenSimulate Strategy //进行仿真时使用的策略出牌策略
}

func NewSimulateStrategy(distCardTimes, simulateTimesPerDist int, simulator Simulator, strategyWhenSimulate Strategy) SimulateStrategy {
	return SimulateStrategy{distCardTimes: distCardTimes, simulateTimesPerDist: simulateTimesPerDist,
		simulator: simulator, strategyWhenSimulate: strategyWhenSimulate}
}

// =================================================
// 仿真策略
// -------------------------------------------------
// 输出的action一定要带Context
func (s SimulateStrategy) MakeDecision(c StrategyContext) (actionTaken action.Action) {
	contextAction := c.PreerAction
	if contextAction.ActionType() == action.Pass {
		contextAction = c.NexterAction
	}
	defer func() { actionTaken = actionTaken.SetContext(contextAction) }()
	possibleActionList := c.RemainPokerSet.PossibleActionsWithContext(contextAction)
	if len(possibleActionList) == 1 {
		return possibleActionList[0]
	}

	//每个action最终有一个总收益
	//key为action 在 possibleActionList中的序号
	actionProfitDict := make([]int, len(possibleActionList))
	partners := [2]int{}
	if c.CheatFlag == WithNexter {
		partners = [2]int{0, 1}
	} else if c.CheatFlag == WithPreer {
		partners = [2]int{0, 1}
	}
	simulator := s.simulator.BuildGame(GameBuildInfo{
		Strategys:   [3]Strategy{s.strategyWhenSimulate, s.strategyWhenSimulate, s.strategyWhenSimulate},
		OpenCard:    c.GameOpenCard,
		Partners:    partners,
		CheatMethod: c.CheatMethod,
	})
	for distIdx := 0; distIdx < s.distCardTimes; distIdx++ {
		if !c.OpenCard4Me { //如果看不到另外两家的牌
			c.PreerRemainPokerSet, c.NexterRemainPokerSet = distCards2(c.OthersRemainPokerSet, c.PreerPkCnt, c.NexterPkCnt, c.PreerPlayLog, c.NexterPlayLog)
		}
		for smIdx := 0; smIdx < s.simulateTimesPerDist; smIdx++ {
			for actionId, a := range possibleActionList {
				a.SetContext(contextAction)
				//simulator中当前player为0号
				simulator.RestoreGame(GameRestoreInfo{
					LastPlayer: 0,
					RemainPokerSets: [3]action.PokerSet2{
						c.RemainPokerSet.Subtract(a.PokerSet2()),
						c.NexterRemainPokerSet,
						c.PreerRemainPokerSet},
					PlayLogs: [3][]action.Action{
						append(c.PlayLog, a),
						c.NexterPlayLog,
						c.PreerPlayLog},
				})
				simulator.RunRestoredGame()
				profits := simulator.GetProfits()
				actionProfitDict[actionId] += profits[0]
				if c.CheatFlag == WithNexter { //with nexter
					actionProfitDict[actionId] += profits[1]
				} else if c.CheatFlag == WithPreer { //with preer
					actionProfitDict[actionId] += profits[2]
				}
			}
		}
	}
	//找出总收益最大的那个action
	return possibleActionList[getMaxProfitIdx(actionProfitDict)]
}

func getMaxProfitIdx(profitDict []int) int {
	maxProfitActionId := 0
	maxProfit := profitDict[0]
	for actionId, profit := range profitDict {
		if profit > maxProfit {
			maxProfit = profit
			maxProfitActionId = actionId
		}
	}
	return maxProfitActionId
}
