package strategy

import (
	"math/rand"

	"github.com/conggova/poker-robot/pkg/action"
)

// 完全随机出牌，不考虑明牌和作弊信息
type RandomStrategy struct {
}

func NewRandomStrategy() RandomStrategy {
	return RandomStrategy{}
}

// =================================================
// 随机策略出牌 ， 在不放走的情况下
// -------------------------------------------------
// 如果下家是一张，要考虑不放走，既保证在随机牌局中不存在放走包赔
func (s RandomStrategy) MakeDecision(c StrategyContext) (actionTaken action.Action) {
	contextAction := c.PreerAction
	if contextAction.ActionType() == action.Pass {
		contextAction = c.NexterAction
	}
	defer func() { actionTaken.SetContext(contextAction) }()
	possibleActionList := c.RemainPokerSet.PossibleActionsWithContext(contextAction)
	if len(possibleActionList) == 1 {
		return possibleActionList[0]
	}

	//如果下家只剩一张，保证不放走包赔
	//如果有可能放走
	//可能放走的条件  1.下家剩一张  2.背景牌没有 或者是单牌
	if c.NexterPkCnt == 1 && (contextAction.ActionType() == action.Pass ||
		contextAction.ActionType() == action.Single) {
		//找有没有非单牌
		var noSingleActionList []action.Action
		for _, a := range possibleActionList {
			if a.ActionType() != action.Single {
				noSingleActionList = append(noSingleActionList, a)
			}
		}
		if len(noSingleActionList) == 0 {
			//找最大的单牌
			maxS := c.RemainPokerSet.MaxPowerPoker()
			return action.NewActionWithoutAff(action.Single, 1, uint64(maxS))
		} else {
			//随机选一个非单牌
			randomIdx := rand.Intn(len(noSingleActionList))
			return noSingleActionList[randomIdx]
		}
	} else {
		//如果不可能放走
		randomIdx := rand.Intn(len(possibleActionList))
		return possibleActionList[randomIdx]
	}
}

// 完全随机出牌，不考虑明牌和作弊信息
// 随机选Times次，最终选张数最多的
type RandomStrategy2 struct {
	times int
}

func NewRandomStrategy2(times int) RandomStrategy2 {
	return RandomStrategy2{times: times}
}

// =================================================
// 随机策略出牌 ，
func (s RandomStrategy2) MakeDecision(c StrategyContext) (actionTaken action.Action) {
	contextAction := c.PreerAction
	if contextAction.ActionType() == action.Pass {
		contextAction = c.NexterAction
	}
	defer func() { actionTaken = actionTaken.SetContext(contextAction) }()
	possibleActionList := c.RemainPokerSet.PossibleActionsWithContext(contextAction)
	if len(possibleActionList) == 1 {
		return possibleActionList[0]
	}

	//如果下家只剩一张，保证不放走包赔
	//如果有可能放走
	//可能放走的条件  1.下家剩一张  2.背景牌没有 或者是单牌
	if c.NexterPkCnt == 1 && (contextAction.ActionType() == action.Pass ||
		contextAction.ActionType() == action.Single) {
		//找有没有非单牌
		var noSingleActionList []action.Action
		for _, a := range possibleActionList {
			if a.ActionType() != action.Single {
				noSingleActionList = append(noSingleActionList, a)
			}
		}
		if len(noSingleActionList) == 0 {
			//找最大的单牌
			maxS := c.RemainPokerSet.MaxPowerPoker()
			return action.NewActionWithoutAff(action.Single, 1, uint64(maxS))
		} else {
			//可以出任意非单牌
			res := noSingleActionList[rand.Intn(len(noSingleActionList))]
			for i := 1; i < s.times; i++ {
				t := noSingleActionList[rand.Intn(len(noSingleActionList))]
				if t.PokerCount() > res.PokerCount() {
					res = t
				}
			}
			return res
		}
	} else if contextAction.ActionType() == action.Pass { //可以出任意牌
		res := possibleActionList[rand.Intn(len(possibleActionList))]
		for i := 1; i < s.times; i++ {
			t := possibleActionList[rand.Intn(len(possibleActionList))]
			if t.PokerCount() > res.PokerCount() {
				res = t
			}
		}
		return res
	} else {
		return possibleActionList[rand.Intn(len(possibleActionList))]
	}
}
