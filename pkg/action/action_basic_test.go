package action

import "testing"

func TestNewActionWithoutAff(t *testing.T) {
	a := NewActionWithoutAff(Single, 5, 5)
	if a[1] != 0x1111100000 || a[0] != 0x05000501 {
		t.Errorf("a should be 0x05000501 0x1111100000 , whereas is %x %x", a[0], a[1])
	}
}

func TestNewActionWithPokerSet2(t *testing.T) {
	a := NewActionWithPokerSet2(SL, 7, 0, 0, 0x1111111)
	if a[1] != 0x1111111 || a[0] != 0x00000706 {
		t.Errorf("NewActionWithPokerSet2 incorrect")
	}
}

func TestActionAttrs(t *testing.T) {
	a := NewActionWithoutAff(Single, 1, 3)
	if a.ActionStructure() != 0x000101 {
		t.Error("ActionStructure incorrect")
	}
	if a.ActionType() != Single {
		t.Error("ActionType incorrect")
	}
	if a.AffLen() != 0 {
		t.Error("AffLen incorrect")
	}
	if a.KeyCard() != 3 {
		t.Error("KeyCard incorrect")
	}
	if a.PokerSet2() != 0x1000 {
		t.Error("PokerSet2 incorrect")
	}
	a = NewActionWithoutAff(Bomb, 1, 5)
	if a.KeyCard() != 5 {
		t.Error("KeyCard incorrect")
	}
	if !a.IsBomb() {
		t.Error("IsBomb incorrect")
	}
	if a.PokerCount() != 4 {
		t.Error("PokerCount incorrect")
	}
	a = a.SetAff(3, 0x12)
	if a.AffLen() != 3 {
		t.Error("SetAff AffLen incorrect")
	}
	if a.ActionType() != BW {
		t.Error("SetAff incorrect")
	}
	if a.PokerCount() != 7 {
		t.Error("SetAff PokerCount incorrect")
	}
	if a.PokerSet2() != 0x400012 {
		t.Error("SetAff incorrect")
	}
	a = a.SetContext(NewActionWithoutAff(Bomb, 1, 4).SetAff(3, 0x12))
	if a.ContextActionType() != BW {
		t.Errorf("SetContext ContextActionType incorrect , %x", a.ContextActionType())
	}
	if a.ContextActionBrief() != 0x04030105 {
		t.Error("ContextActionBrief incorrect")
	}
	if a.WithoutContext() != 0x05030105 {
		t.Error("WithoutContext incorrect")
	}
}

func TestPlayableInContext(t *testing.T) {
	a := NewActionWithoutAff(Bomb, 1, 4)
	c := a
	if a.PlayableInContext(c) {
		t.Error("PlayableInContext incorrect1")
	}
	c = NewActionWithoutAff(Bomb, 1, 3)
	if !a.PlayableInContext(c) {
		t.Error("PlayableInContext incorrect2")
	}
	a = NewActionWithoutAff(SL, 6, 3)
	c = NewActionWithoutAff(SL, 5, 3)
	if a.PlayableInContext(c) {
		t.Error("PlayableInContext incorrect3")
	}
	a.SetContext(c)
	c = NewActionWithoutAff(SL, 6, 2)
	c.SetContext(a)
	if !a.PlayableInContext(c) {
		t.Error("PlayableInContext incorrect4")
	}
	a = NewActionWithoutAff(Bomb, 1, 4)
	if !a.PlayableInContext(c) {
		t.Error("PlayableInContext incorrect5")
	}
}
