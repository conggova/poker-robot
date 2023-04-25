package player

import (
	"fmt"

	"github.com/conggova/poker-robot/pkg/action"
	"github.com/conggova/poker-robot/pkg/strategy"
)

// 人类玩家,不参与作弊
type ManualPlayer struct {
	BasePlayer
}

func NewManualPlayer(playerId int) *ManualPlayer {
	return &ManualPlayer{BasePlayer{Id: playerId}}
}

func printPlayersInfo(opencard bool, remainPokerSet, preerPokerSet, nexterPokerSet action.PokerSet2,
	preerCnt, nexterCnt int, preerAction, nexterAction action.Action) {
	fmt.Println("")
	fmt.Println("当前情况:")
	fmt.Println("--------------------------------------------------------")
	fmt.Println("你的上家出 ", preerAction)
	if opencard {
		fmt.Println("余牌：", preerPokerSet)
	} else {
		fmt.Printf("当前还剩 %d 张", preerCnt)
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("                         你的下家出 ", nexterAction)
	if opencard {
		fmt.Println("                         余牌：", nexterPokerSet)
	} else {
		fmt.Printf("                         当前还剩 %d 张", nexterCnt)
	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("你现在手里有：", remainPokerSet)
	fmt.Println("该你出牌")
	fmt.Println("--------------------------------------------------------")
	fmt.Println("")
}

func (p *ManualPlayer) Play(pc strategy.PlayContext) (actionTaken action.Action) {
	var contextAction = pc.PreerAction
	if contextAction.ActionType() == action.Pass {
		contextAction = pc.NexterAction
	}
	defer func() {
		actionTaken = actionTaken.SetContext(contextAction)
		p.remainPokerSet = p.remainPokerSet.Subtract(actionTaken.PokerSet2())
	}()
	//有走必走
	possibleActionList := p.remainPokerSet.PossibleActionsWithContext(contextAction)
	//如果只有一种选择
	if len(possibleActionList) == 1 {
		possibleActionList[0] = possibleActionList[0].SetContext(contextAction)
		printPlayersInfo(pc.GameOpenCard, p.remainPokerSet, pc.PreerRemainPokerSet, pc.NexterRemainPokerSet, pc.PreerPkCnt, pc.NexterPkCnt, pc.PreerAction, pc.NexterAction)
		fmt.Println("Your only choice is : ", possibleActionList[0])
		fmt.Print("Type any key to continue......")
		var appLine string
		fmt.Scanln(&appLine)
		return possibleActionList[0]
	}

	for {
		printPlayersInfo(pc.GameOpenCard, p.remainPokerSet, pc.PreerRemainPokerSet, pc.NexterRemainPokerSet, pc.PreerPkCnt, pc.NexterPkCnt, pc.PreerAction, pc.NexterAction)
		var ipt string
		fmt.Print("Input (type H for help):")
		fmt.Scanln(&ipt)

		//打印帮助
		if ipt == "H" {
			printHelpInfo()
			continue
		}

		//检查输入
		ok, checkResult := checkInput(ipt, p.remainPokerSet, contextAction)
		//如果是过 或者输入有效
		if ok {
			fmt.Println("Input is :", checkResult)
			var ynInput string
			fmt.Print("Type R for retry , others for confirm :")
			fmt.Scanln(&ynInput)
			confirm := true
			if ynInput == "R" {
				confirm = false
			}

			//只有输入正确 ， 并且确认的情况下
			if confirm {
				return checkResult
			} else {
				fmt.Println("")
				fmt.Println("Discarded last input , retry.")
			}

		} else { //输入不符合规则
			fmt.Println("Input is not proper. Please Retry!!!")
			fmt.Println("")
		}
	}
}

func checkInput(ipt string, ps action.PokerSet2, contextAction action.Action) (bool, action.Action) {
	//先检查输入是否符合规则
	ok, a := action.ParseAction(ipt)
	if !ok {
		return false, action.Action{}
	}
	//再检查输入是否符合背景
	if !a.PlayableInContext(contextAction) {
		return false, action.Action{}
	}
	//再检查是否能提供
	if !a.CanBeSupliedBy(ps) {
		return false, action.Action{}
	}
	return true, a
}

func printHelpInfo() {
	fmt.Println("#--------------------------------------#")
	fmt.Println("HELP INFO ：")
	fmt.Println("每张牌的表示：3456789TJQKA2，特别注意T代表10")
	fmt.Println("各种牌形的介绍如下：")
	fmt.Println("单牌: J(一个J)")
	fmt.Println("对子: TT(对10)")
	fmt.Println("三条: 999(三个9)")
	fmt.Println("三带一: 999K(三个9带K)")
	fmt.Println("三带二: QQQ34(三个Q带34)")
	fmt.Println("炸弹: KKKK")
	fmt.Println("四带一: 77778(四个7带8)")
	fmt.Println("四带二: 777789(四个7带8和9)")
	fmt.Println("四带三: 777789T(四个7带89T)")
	fmt.Println("顺子: 3456789TJQKA(顺子3到A)")
	fmt.Println("双顺: 6677(双顺6到7)")
	fmt.Println("飞机不带: 333444(333444)")
	fmt.Println("飞机带一: 33344456(333444带56)")
	fmt.Println("飞机带二: 3334445567(333444带5567)")
	fmt.Println("具体规则可以参照：https://jingyan.baidu.com/article/e9fb46e1c7becb3521f7668f.html")
	fmt.Println("#--------------------------------------#")
}
