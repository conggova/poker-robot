package strategy

import (
	"fmt"
	"math/rand"

	"github.com/conggova/poker-robot/pkg/action"
)

/*
var idxPkNumMap []int = []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
	3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
*/
/*
同时返回哪一位是头家
发牌函数,发32张牌大概需要51次
除了第一次，其余每次发牌所需要的次数是无穷等比数列
1+（i-1）/48 + (（i-1）/48)**2 ...
总数计算：
sum = 0
for i in range(0,32):

	sum += 48/(48-i)

print(sum)
*/
func DistCards() [3]action.PokerSet {
	var res [3]action.PokerSet
	//1代表已被选过了
	var choosedMark uint64
	//生成前两手
	pkCnt := 0
	for {
		randomIdx := uint64(1 << rand.Intn(48))
		if choosedMark&randomIdx == 0 { //not choosen yet
			choosedMark |= randomIdx
			pkCnt += 1
			if pkCnt == 16 {
				break
			}
		}
	}
	res[0] = action.PokerSet(choosedMark)

	pkCnt = 0
	for {
		randomIdx := uint64(1 << rand.Intn(48))
		if choosedMark&randomIdx == 0 { //not choosen yet
			choosedMark |= randomIdx
			pkCnt += 1
			if pkCnt == 16 {
				break
			}
		}
	}
	res[1] = action.PokerSet(choosedMark).Subtract(res[0])
	res[2] = action.PokerSet(0b111111111111111111111111111111111111111111111111).Subtract(action.PokerSet(choosedMark))
	for i := 0; i < 3; i++ {
		if res[i]&0b100000000000000000000000000000000000000000000000 != 0 { //move 48 to 49
			res[i] &= 0b011111111111111111111111111111111111111111111111
			res[i] |= 0b1000000000000000000000000000000000000000000000000
			break
		}
	}
	return res
}

// 切牌函数
func distCards2(ps action.PokerSet2, preerNum int,
	nexterNum int, preerPlayLog []action.Action, nexterPlayLog []action.Action) (action.PokerSet2, action.PokerSet2) {
	//当初过的牌，现在必须能过, 新发的牌+patch必不能Afford以前Pass的牌
	preerPatchActionsMap := map[action.PokerSet2][]uint64{}
	var patch action.PokerSet2
	for idx := len(preerPlayLog) - 1; idx >= 0; idx-- {
		a := preerPlayLog[idx]
		if a.ActionType() == action.Pass && a.ContextActionType() != action.Pass {
			preerPatchActionsMap[patch] = append(preerPatchActionsMap[patch], a.ContextActionBrief())
		} else {
			patch = patch.CombineWith(a.PokerSet2())
		}
	}

	nexterPatchActionsMap := map[action.PokerSet2][]uint64{}
	patch = 0
	for idx := len(nexterPlayLog) - 1; idx >= 0; idx-- {
		a := nexterPlayLog[idx]
		if a.ActionType() == action.Pass && a.ContextActionType() != action.Pass {
			nexterPatchActionsMap[patch] = append(nexterPatchActionsMap[patch], a.ContextActionBrief())
		} else {
			patch = patch.CombineWith(a.PokerSet2())
		}
	}
	trycnt := 0
	for {
		trycnt += 1
		/*
			if trycnt > 1000 {
				fmt.Println("to divide : ", ps, " preerNum , nexterNum : ", preerNum, nexterNum)
				fmt.Println(preerPlayLog)
				fmt.Println(nexterPlayLog)
				fmt.Println(preerPatchActionsMap)
				fmt.Println(nexterPatchActionsMap)
			}
		*/
		if trycnt%1000 == 0 {
			fmt.Println("cut card on going , please wait. Already tried ", trycnt, " times.")
		}
		prePS, nextPS := cutCardsRandom(ps, preerNum, nexterNum)
		if checkConsistantness3(preerPatchActionsMap, nexterPatchActionsMap, prePS, nextPS) {
			//fmt.Println(trycnt)
			return prePS, nextPS
		}
	}

}

func checkConsistantness3(preerPatchActionsMap, nexterPatchActionsMap map[action.PokerSet2][]uint64, prePS, nextPS action.PokerSet2) bool {
	for patch, actions := range preerPatchActionsMap {
		ps := patch.CombineWith(prePS)
		for _, a := range actions {
			if ps.Afford(a) {
				return false
			}
		}
	}
	for patch, actions := range nexterPatchActionsMap {
		ps := patch.CombineWith(nextPS)
		for _, a := range actions {
			if ps.Afford(a) {
				return false
			}
		}
	}
	return true
}

func cutCardsRandom(ps action.PokerSet2, preerNum int, nexterNum int) (action.PokerSet2, action.PokerSet2) {
	cards := []int8{} //每个元素一张牌
	for i := 0; i < 15; i++ {
		cnt := ps >> (i << 2) & 0b1111
		for j := 0; j < int(cnt); j++ {
			cards = append(cards, int8(i))
		}
	}
	if preerNum+nexterNum != len(cards) {
		panic(fmt.Sprintf("preerNum %d + nexterNum %d != len(cards) %d , %s ", preerNum, nexterNum, len(cards), ps))
	}
	//true代表已被选过了
	choosed := make([]bool, len(cards))
	var prePS action.PokerSet2
	pkCnt := 0
	for {
		randomIdx := rand.Intn(int(preerNum + nexterNum))
		if !choosed[randomIdx] {
			choosed[randomIdx] = true
			prePS += 1 << (cards[randomIdx] << 2)
			pkCnt += 1
			if pkCnt == preerNum {
				break
			}
		}
	}
	return prePS, ps.Subtract(prePS)
}

/*
func checkConsistantness(prePS, nextPS action.PokerSet,
	preerPlayLog []action.Action, nexterPlayLog []action.Action) bool {
	preerPassedActions := []action.Action{}
	for _, a := range preerPlayLog {
		if a.ActionType == action.Pass && a.Context.ActionType != action.Pass {
			preerPassedActions = append(preerPassedActions, *a.Context)
		}
	}
	nexterPassedActions := []action.Action{}
	for _, a := range nexterPlayLog {
		if a.ActionType == action.Pass && a.Context.ActionType != action.Pass {
			nexterPassedActions = append(nexterPassedActions, *a.Context)
		}
	}
	//当初过的牌，现在必须能过
	//还应该有更严格的条件  还原到过去要不起的时刻
	for _, a := range preerPassedActions {
		//如果现在要得起，说明不一致
		possibleActions := action.FindAllActions(prePS, a)
		//如果action里面不是只有过牌一项 ， 说明不一致
		if !(len(possibleActions) == 1 && possibleActions[0].ActionType == action.Pass) {
			return false
		}
	}
	//当初过的牌，现在必须能过
	for _, a := range nexterPassedActions {
		//如果现在要得起，说明不一致
		possibleActions := action.FindAllActions(nextPS, a)
		//如果action里面不是只有过牌一项 ， 说明不一致
		if !(len(possibleActions) == 1 && possibleActions[0].ActionType == action.Pass) {
			return false
		}
	}
	return true
}
*/
