package strategy

import (
	"testing"

	"github.com/conggova/poker-robot/pkg/action"
)

func Test_distCards2(t *testing.T) {
	p1, p2 := distCards2(0x102200000, 1, 4,
		[]action.Action{
			action.NewActionWithoutAff(action.Pass, 0, 0).SetContext(action.NewActionWithoutAff(action.Triple, 1, 4)),
			action.NewActionWithoutAff(action.CL, 2, 5)},
		[]action.Action{})
	if p1.CombineWith(p2) != 0x102200000 || p1.PokerCount() != 1 || p2.PokerCount() != 4 {
		t.Error("distCards2 incorrect")
	}
	if p1 != 0x100000000 {
		t.Error("distCards2 incorrect2", p1, p2)
	}
}
