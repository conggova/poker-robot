package game

import (
	"fmt"
	"strconv"

	"github.com/conggova/poker-robot/pkg/player"
	"github.com/conggova/poker-robot/pkg/strategy"
)

// 人机对战
func ManAndRobot() {
	var ipt string
	gameOpenCard := false
	fmt.Println("您希望游戏是明牌吗？（大家都能看到彼此的牌）")
	fmt.Print("Y for yes , others for no :")
	ipt = "N"
	fmt.Scanln(&ipt)
	if ipt == "Y" {
		gameOpenCard = true
	}

	othersAreParterners := false
	othersCom := false
	fmt.Println("您希望另外两个玩家是同伙吗？（利益共同体）")
	fmt.Print("Y for yes , others for no :")
	ipt = "N"
	fmt.Scanln(&ipt)
	if ipt == "Y" {
		othersAreParterners = true
		if !gameOpenCard {
			fmt.Println("您希望他们能够私下通信吗?（能相互看牌）")
			fmt.Print("Y for yes , others for no :")
			ipt = "N"
			fmt.Scanln(&ipt)
			if ipt == "Y" {
				othersCom = true
			}
		}
	}

	gameTotal := 0
	fmt.Print("您想玩几局？")
	fmt.Scanf("%d\n", &gameTotal)

	manTotalProfit := 0
	for i := 0; i < gameTotal; i++ {
		fmt.Println("一局游戏开始 ---------------------------------------------")
		manProfit := manAndRobotPlayOneGame(gameOpenCard, othersAreParterners, othersCom)
		fmt.Println("这一局您的收益是 ", manProfit)
		manTotalProfit += manProfit
		fmt.Println("目前您的总收益是 ", manTotalProfit)
		fmt.Println("")
		fmt.Println("此局结束-------------------------------------------------")
	}
	fmt.Println("你一共玩了 ", gameTotal, " 局牌， 总收益是 ", manTotalProfit, " 。 ")
	fmt.Println("")
	fmt.Println("")

	fmt.Println("您要再玩几局吗？")
	fmt.Print("Y for yes , others for no :")
	ipt = "N"
	fmt.Scanln(&ipt)
	if ipt == "Y" {
		ManAndRobot()
	}
}

func manAndRobotPlayOneGame(gameOpenCard bool, robotsArePartners bool, robotsCom bool) int {
	var player1 player.Player = player.NewManualPlayer(0)
	var player2 player.Player = player.NewStrategyPlayer(1, strategy.NewSimulateStrategy(50, 5, &Game{}, strategy.NewRandomStrategy2(10)), strategy.NoCheat, strategy.ShareInterest, nil)
	var player3 player.Player = player.NewStrategyPlayer(2, strategy.NewSimulateStrategy(50, 5, &Game{}, strategy.NewRandomStrategy2(10)), strategy.NoCheat, strategy.ShareInterest, nil)

	if robotsArePartners {
		player2.SetCheatFlag(strategy.WithNexter)
		player3.SetCheatFlag(strategy.WithPreer)
		if robotsCom {
			player2.SetCheatMethod(strategy.CommInSecret)
			player3.SetCheatMethod(strategy.CommInSecret)
		}
		player2.SetPartner(player3)
		player3.SetPartner(player2)
	}

	game := newGame([3]player.Player{player1, player2, player3}, gameOpenCard, true)
	game.RunGameWithRandomBeginning()
	manProfit := game.GetProfits()[0]
	return manProfit
}

// 机器人大战
func RobotFight() {
	totalProfits := [3]int{}
	for i := 1; i < 10000; i++ {
		fmt.Println("第 ", i, " 局")
		profits := robotFightPlayOneGame(false)
		for i := 0; i < 3; i++ {
			totalProfits[i] += profits[i]
		}
		fmt.Print("\n")
		fmt.Print("总收益：\n")
		fmt.Print("	刘备 ")
		fmt.Print(strconv.Itoa(totalProfits[0]))
		fmt.Print("\n")
		fmt.Print("	关羽 ")
		fmt.Print(strconv.Itoa(totalProfits[1]))
		fmt.Print("\n")
		fmt.Print("	张飞 ")
		fmt.Print(strconv.Itoa(totalProfits[2]))
		fmt.Print("\n")
		fmt.Print("\n")
		fmt.Print("回车继续... ")
		var ipt string
		fmt.Scanln(&ipt)
	}
}

func robotFightPlayOneGame(opencard bool) [3]int {
	var player1 player.Player = player.NewStrategyPlayer(0, strategy.NewSimulateStrategy(5, 10, &Game{}, strategy.NewRandomStrategy()), strategy.NoCheat, strategy.ShareInterest, nil)
	var player2 player.Player = player.NewStrategyPlayer(1, strategy.NewSimulateStrategy(5, 10, &Game{}, strategy.NewRandomStrategy()), strategy.NoCheat, strategy.ShareInterest, nil)
	var player3 player.Player = player.NewStrategyPlayer(2, strategy.NewSimulateStrategy(5, 10, &Game{}, strategy.NewRandomStrategy2(10)), strategy.NoCheat, strategy.ShareInterest, nil)
	game := newGame([3]player.Player{player1, player2, player3}, opencard, true)
	game.RunGameWithRandomBeginning()
	return game.GetProfits()
}
