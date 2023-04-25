package action

import (
	"testing"
)

func TestCanBeSupliedBy(t *testing.T) {
	a := NewActionWithoutAff(SL, 5, 3)
	if !a.CanBeSupliedBy(0x22222000) {
		t.Error("CanBeSupliedBy incorrect")
	}
	if a.CanBeSupliedBy(0x11111) {
		t.Error("CanBeSupliedBy incorrect")
	}
}

func TestParseAction(t *testing.T) {
	if ok, a := ParseAction("33345"); !ok || a.ActionType() != Triple || a.AffLen() != 2 {
		t.Error("ParseAction incorrect")
	}

	if ok, a := ParseAction("3344"); !ok || a.ActionType() != CL || a.AffLen() != 0 || a.SeqLen() != 2 {
		t.Error("ParseAction incorrect")
	}

	if ok, a := ParseAction("3334445568"); !ok || a.ActionType() != TL || a.AffLen() != 4 || a.SeqLen() != 2 {
		t.Error("ParseAction incorrect")
	}

	if ok, a := ParseAction("33344455689"); ok || a.ActionType() == TL {
		t.Error("ParseAction incorrect")
	}

	if ok, a := ParseAction("444489"); !ok || a.ActionType() != BW || a.AffLen() != 2 || a.SeqLen() != 1 {
		t.Error("ParseAction incorrect")
	}
}

func Test_getSubSeqsWithMinLen(t *testing.T) {
	seqs := getSubSeqsWithMinLen([2]uint64{2, 7}, 5)
	if len(seqs) != 6 {
		t.Error("getSubSeqsWithMinLen incorrect")
	}
}

func Test_getSubSeqsWithFixLen(t *testing.T) {
	seqs := getSubSeqsWithFixLen([2]uint64{2, 7}, 5)
	if len(seqs) != 3 || seqs[0] != [2]uint64{2, 5} {
		t.Error("getSubSeqsWithFixLen incorrect")
	}
}

func Test_getAffs(t *testing.T) {
	affs := getAffs(0x321, 0, 4)
	if len(affs) != 5 {
		t.Error("getAffs incorrect")
	}
	affs = getAffs(0x21200100010, 0, 3)
	if len(affs) != 18 {
		t.Error("getAffs incorrect")
	}

	affs = getAffs(0x21200100010, 0, 8)
	if len(affs) != 18 {
		t.Error("getAffs incorrect")
	}

}

func Test_possibleActionsWithoutContext(t *testing.T) {
	var ps PokerSet2 = 0x21
	actions := ps.possibleActionsWithoutContext()
	if len(actions) != 3 {
		t.Error("possibleActionsWithoutContext incorrect")
	}
	ps = 0x4
	actions = ps.possibleActionsWithoutContext()
	if len(actions) != 1 {
		t.Error("possibleActionsWithoutContext incorrect")
	}

	ps = 0x111110
	actions = ps.possibleActionsWithoutContext()
	if len(actions) != 1 {
		t.Error("possibleActionsWithoutContext incorrect")
	}

	ps = 0x2220
	actions = ps.possibleActionsWithoutContext()
	if len(actions) != 1 {
		t.Error("possibleActionsWithoutContext incorrect")
	}

	ps = 0x33310
	actions = ps.possibleActionsWithoutContext()
	if len(actions) != 1 {
		t.Error("possibleActionsWithoutContext incorrect")
	}

	ps = 0x421
	actions = ps.possibleActionsWithoutContext()
	if len(actions) == 1 {
		t.Error("possibleActionsWithoutContext incorrect")
	}
}

func Test_possibleActionsWithContext(t *testing.T) {
	var ps PokerSet2 = 0x21
	actions := ps.possibleActionsWithContext(NewActionWithoutAff(Couple, 1, 5))
	if len(actions) != 0 {
		t.Error("possibleActionsWithoutContext incorrect")
	}
	ps = 0x2000001
	actions = ps.possibleActionsWithContext(NewActionWithoutAff(Couple, 1, 5))
	if len(actions) != 1 {
		t.Error("possibleActionsWithoutContext incorrect")
	}
}

func TestPossibleActionsWithContext(t *testing.T) {
	var ps PokerSet2 = 0x21
	actions := ps.PossibleActionsWithContext(NewActionWithoutAff(Couple, 1, 5))
	if len(actions) != 1 || actions[0].ActionType() != Pass {
		t.Error("PossibleActionsWithContext incorrect")
	}
}

func TestAfford(t *testing.T) {
	var ps PokerSet2 = 0x211133000
	if ps.Afford(NewActionWithoutAff(SL, 6, 3)[0]) {
		t.Error("Afford incorrect")
	}
	ps = 0x1211133000
	if !ps.Afford(NewActionWithoutAff(SL, 6, 3)[0]) {
		t.Error("Afford incorrect")
	}

	ps = 0x21043033
	if ps.Afford(NewActionWithoutAff(Bomb, 1, 10)[0]) {
		t.Error("Afford incorrect2")
	}
}
