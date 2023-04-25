package game

import (
	"fmt"
	"os"
	"strconv"

	"github.com/conggova/poker-robot/pkg/action"
)

type GameStat struct {
	keepLog           bool                //如果KeepLog会写日志到文件
	playerRemainCards [3]action.PokerSet2 //最后剩下的牌
	playLog           []action.Action     //行牌日志
	lastPlayer        int                 //最后的出牌人，牌局结束时即为赢家
	bomberLog         []int               //成立的炸弹,元素为0 1 2代表对应的玩家
	letGoer           int                 //放走者
	profits           [3]int              //记录收益
	end               bool                //游戏已经结束
	analyzed          bool                //已经分析过
}

func newGameStat(keepLog bool) GameStat {
	return GameStat{keepLog: keepLog}
}

// 记录一次出牌，playLog有可能是不全的
func (gs *GameStat) putAction(playerId int, a action.Action) {
	l := len(gs.playLog)
	gs.lastPlayer = playerId
	if l >= 2 {
		//处理炸弹,如果被反炸，则炸弹不成立
		//Bomb接两手Pass，则炸弹成立
		if a.ActionType() == action.Pass && gs.playLog[l-1].ActionType() == action.Pass && gs.playLog[l-2].IsBomb() {
			gs.bomberLog = append(gs.bomberLog, (playerId+1)%3)
		}
		if gs.keepLog { //持续追加
			gs.playLog = append(gs.playLog, a)
		} else { //移位，保持两条
			gs.playLog[0] = gs.playLog[1]
			gs.playLog[1] = a
		}
	} else {
		gs.playLog = append(gs.playLog, a)
	}
}

// 玩家亮底牌，游戏结束,确定赢家
func (gs *GameStat) endGame(playerRemainCards [3]action.PokerSet2) {
	gs.end = true
	gs.playerRemainCards = playerRemainCards
	//记录最后一个炸弹
	if gs.playLog[len(gs.playLog)-1].IsBomb() {
		gs.bomberLog = append(gs.bomberLog, gs.lastPlayer)
	}
}

func (gs *GameStat) analyze() {
	if !gs.end {
		panic("must end game before u can analyze")
	}

	gs.analyzed = true
	gs.calcuProfit()
	if gs.keepLog {
		gs.logGame(true)
	}
}

func (gs *GameStat) getProfits() [3]int {
	if !gs.analyzed {
		panic("must analyze before getProfits")
	}
	return gs.profits
}

// 计算每个玩家的收益
// 考虑的因素
// 是否被关  是否有成立的炸弹  是否放走
func (gs *GameStat) calcuProfit() {
	//炸弹
	for _, playerId := range gs.bomberLog {
		gs.profits[playerId] += 2 * 5 * 5
		gs.profits[(playerId+2)%3] -= 5 * 5
		gs.profits[(playerId+1)%3] -= 5 * 5
	}
	//赢家的上下家
	winnerPreId := (gs.lastPlayer + 2) % 3
	winnerNextId := (gs.lastPlayer + 1) % 3
	//包赔判断是否存在包赔
	//最后一手赢家是单，赢家上家最后一手是单
	haveLetgoer := false
	l := len(gs.playLog)
	//赢家最后出的是单
	if gs.playLog[l-1].ActionType() == action.Single {
		winnerPreLastHand := gs.playLog[l-2]
		//赢家前一家最后出的是单
		if winnerPreLastHand.ActionType() == action.Single {
			//如果赢上家出牌的背景是Pass
			if winnerPreLastHand.ContextActionType() == action.Pass {
				//还原其出牌时的牌，看看有没有可能不出单 ， 如果可能 ， 那么就是放走
				//如果只有单牌 再看出的是不是最大的单牌
				if gs.playerRemainCards[winnerPreId].MaxPowerPoker() > int(winnerPreLastHand.KeyCard()) ||
					gs.playerRemainCards[winnerPreId].CombineWith(winnerPreLastHand.PokerSet2()).HaveNonSingle() {
					haveLetgoer = true
				}
				/*
					if action.GetMaxSingleNum(gs.playerRemainCards[winnerPreId]) > winnerPreLastHand.Main1 ||
						action.HaveNonSingle(action.CombinedPkStatMap(gs.playerRemainCards[winnerPreId], winnerPreLastHand.PkStatMap())) {
						haveLetgoer = true
					}
				*/
			} else { //出牌背景为单
				//看看他手中有没有更大的单牌
				if gs.playerRemainCards[winnerPreId].MaxPowerPoker() > int(winnerPreLastHand.KeyCard()) {
					haveLetgoer = true
				}
				/*
					if action.GetMaxSingleNum(gs.playerRemainCards[winnerPreId]) > winnerPreLastHand.Main1 {
						haveLetgoer = true
					}
				*/
			}
		}
	}

	winnerPreRemain := gs.playerRemainCards[winnerPreId].PokerCount()
	winnerNextRemain := gs.playerRemainCards[winnerNextId].PokerCount()
	if haveLetgoer {
		gs.letGoer = winnerPreId
		winnerTotalWin := 0
		if winnerPreRemain == 16 {
			winnerTotalWin += 10 * 16
		} else {
			winnerTotalWin += 5 * winnerPreRemain
		}

		if winnerNextRemain == 16 {
			winnerTotalWin += 10 * 16
		} else {
			winnerTotalWin += 5 * int(winnerNextRemain)
		}
		gs.profits[gs.lastPlayer] += winnerTotalWin
		gs.profits[gs.letGoer] -= winnerTotalWin
	} else {
		preLose := 0
		if winnerPreRemain == 16 {
			preLose += 10 * 16
		} else {
			preLose += 5 * int(winnerPreRemain)
		}
		nextLose := 0
		if winnerNextRemain == 16 {
			nextLose += 10 * 16
		} else {
			nextLose += 5 * int(winnerNextRemain)
		}
		gs.profits[gs.lastPlayer] += preLose + nextLose
		gs.profits[winnerPreId] -= preLose
		gs.profits[winnerNextId] -= nextLose
	}
}

// 记录出牌日志
func (gs *GameStat) logGame(verbose bool) {
	headPlayer := (gs.lastPlayer - (len(gs.playLog) - 1) + 3*100) % 3 //可以从Winner推导 , (head + cnt -1)%3 = last
	//还原初始牌
	var playerInitCards [3]action.PokerSet2 = gs.playerRemainCards
	for i, a := range gs.playLog {
		p := (headPlayer + i) % 3
		playerInitCards[p] = playerInitCards[p].CombineWith(a.PokerSet2())
	}

	fout, err := os.OpenFile("gameLog.txt", os.O_APPEND|os.O_CREATE, 0755)
	defer fout.Close()
	fwrite := func(s string) {
		if verbose {
			fmt.Print(s)
		}
		fout.WriteString(s)
	}
	if err != nil {
		fmt.Println("notes.txt", err)
		return
	}
	fwrite("")

	fwrite("----------------牌局记录开始---------------\n")
	currentPlayer := headPlayer
	loopCnt := 0
	nameMap := map[int]string{0: "刘备", 1: "关羽", 2: "张飞"}
	for idx := 0; idx < len(gs.playLog); idx++ {
		pkStatMapStr := playerInitCards[currentPlayer].String()
		playerInitCards[currentPlayer] = playerInitCards[currentPlayer].Subtract(gs.playLog[idx].PokerSet2())
		//action.ExtractActionFromStatMap(&playerInitCards[currentPlayer], gs.playLog[idx])
		name := nameMap[currentPlayer]
		if currentPlayer == headPlayer {
			loopCnt = loopCnt + 1
			fwrite("第")
			fwrite(strconv.Itoa(loopCnt))
			fwrite("圈-----------\n")
			fwrite("\n")
		}

		fwrite(name)
		fwrite(" 出牌... \n")
		fwrite("	他手中有 ")
		fwrite(pkStatMapStr)
		fwrite("\n")
		fwrite("	他出  ")
		fwrite(gs.playLog[idx].String())
		fwrite("\n")
		fwrite("\n")
		currentPlayer = (currentPlayer + 1) % 3
	}
	fwrite("各方收益：\n")
	fwrite("	刘备 ")
	fwrite(strconv.Itoa(gs.profits[0]))
	fwrite("\n")
	fwrite("	关羽 ")
	fwrite(strconv.Itoa(gs.profits[1]))
	fwrite("\n")
	fwrite("	张飞 ")
	fwrite(strconv.Itoa(gs.profits[2]))
	fwrite("\n")
	fwrite("----------------牌局记录结束---------------\n")
	fwrite("\n")
	fwrite("\n")
}
