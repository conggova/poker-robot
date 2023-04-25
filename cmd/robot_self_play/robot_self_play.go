// pdkRobot project main.go
package main

import "github.com/conggova/poker-robot/pkg/game"

//"fmt"
//"os"

func main() {
	game.RobotFight()
	/*
		var ps action.PokerSet2 = 0x1111111111111
		fmt.Println(ps)
		_, a := action.ParseAction("34567")
		fmt.Println(a)
		actions := ps.PossibleActionsWithContext([2]uint64{0, 0})
		for _, a := range actions {
			fmt.Printf("%x , %x , %s \n", a[0], a[1], a)
		}
	*/

}
