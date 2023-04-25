package action

import (
	"testing"
)

func TestPokerSetBasic(t *testing.T) {
	var ps PokerSet = 0b00010000001101011111
	if ps.PokerSet2() != 0x10224 {
		t.Error("PokerSet PokerSet2 incorrect")
	}
	if ps.PokerCount() != 9 {
		t.Errorf("PokerSet PokerCount incorrect , %d", ps.PokerCount())
	}
	if ps.CombineWith(0b11000000000000000000) != 0b11010000001101011111 {
		t.Error("PokerSet CombineWith incorrect")
	}
	ps = ps.CombineWith(0b11000000000000000000)
	if ps.Subtract(0b11000000000000000000) != 0b00010000001101011111 {
		t.Error("PokerSet Subtract incorrect")
	}
	if ps.MaxPowerPoker() != 4 {
		t.Error("PokerSet MaxPowerPoker incorrect")
	}
	if ps.String() != "33334455777" {
		t.Error("PokerSet String incorrect ", ps.String())
	}
}

func TestPokerSet2Basic(t *testing.T) {
	var ps PokerSet2 = 0x10224
	if ps.PokerCount() != 9 {
		t.Errorf("PokerSet2 PokerCount incorrect , %d", ps.PokerCount())
	}
	if ps.CombineWith(0x21000) != 0x31224 {
		t.Error("PokerSet2 CombineWith incorrect")
	}
	ps = ps.CombineWith(0x21000)
	if ps.Subtract(0x21000) != 0x10224 {
		t.Error("PokerSet2 Subtract incorrect")
	}
	if ps.MaxPowerPoker() != 4 {
		t.Error("PokerSet2 MaxPowerPoker incorrect")
	}
	if ps.String() != "333344556777" {
		t.Error("PokerSet2 String incorrect ", ps.String())
	}
	if !ps.HaveNonSingle() {
		t.Error("PokerSet2 HaveNonSingle incorrect1 ")
	}
	ps = 0x1111100
	if !ps.HaveNonSingle() {
		t.Error("PokerSet2 HaveNonSingle incorrect2 ")
	}
	if !ps.Covers(ps) {
		t.Error("PokerSet2 Covers incorrect1")
	}
	if !ps.Covers(0x0111100) {
		t.Error("PokerSet2 Covers incorrect2")
	}
	if ps.Covers(0x1111110) {
		t.Error("PokerSet2 Covers incorrect3")
	}
	if ps.Covers(0x1110) {
		t.Error("PokerSet2 Covers incorrect4")
	}
}

func TestParsePokerSet2(t *testing.T) {
	if ok, v := ParsePokerSet2("33344556tJQ2"); !ok || v != 0x1001110001223 {
		t.Error("PokerSet2 ParsePokerSet2 incorrect")
	}
	if ok, v := ParsePokerSet2("tjq333445562"); !ok || v != 0x1001110001223 {
		t.Error("PokerSet2 ParsePokerSet2 incorrect")
	}
	if ok, _ := ParsePokerSet2("3334455g6tJQ2"); ok {
		t.Error("PokerSet2 ParsePokerSet2 incorrect")
	}
}

func Test_parse(t *testing.T) {
	var ps PokerSet2 = 0x31224
	pkCnt, sCnt, cCnt, tCnt, bCnt, sList, cList, tList, bList := ps.parse()
	if pkCnt != 12 || sCnt != 5 || cCnt != 4 || tCnt != 2 || bCnt != 1 || len(sList) != 1 || len(cList) != 2 || len(tList) != 2 || len(bList) != 1 {
		t.Error("PokerSet2 parse incorrect")
	}
	if sList[0] != [2]uint64{0, 5} || cList[0] != [2]uint64{0, 3} || tList[1] != [2]uint64{4, 1} || bList[0] != 0 {
		t.Error("PokerSet2 parse incorrect")
	}
}
