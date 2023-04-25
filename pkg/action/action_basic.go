// Action的基本定义的操作
package action

// possible types of a action
const (
	Pass   = 0 //过牌
	Single = 1 //单牌
	Couple = 2 //对
	Triple = 3 //三 0 1 2
	Bomb   = 4 //炸
	BW     = 5 //四带 1 2 3
	SL     = 6 //单顺
	CL     = 7 //双顺 2起步
	TL     = 8 //三顺
)

var actionTypeNameMap = map[uint64]string{
	0: "过牌",
	1: "单牌",
	2: "对子",
	3: "三带",
	4: "炸弹",
	5: "四带",
	6: "顺子",
	7: "连对",
	8: "飞机",
}

var actionTypeBaseLenMap = map[uint64]int{
	0: 0,
	1: 1,
	2: 2,
	3: 3,
	4: 4,
	5: 4,
	6: 1,
	7: 2,
	8: 3,
}

// for first uint64
// 8 bits for action type , 8 bit for seqlen , 8 bits for afflen , 8 bits for keycard
// 32 bits for context action , struct is the same as above
//
//	|----------context-------------| |--keycard--| |- afflen -| |- seqlen -| |- actionType -|
//
// 0b00000000000000000000000000000000   00000000      00000000     00000000      00000000
// for second uint64 represent a PokerSet2
type Action [2]uint64

func genPokerSet2WithActionStructionInfo(actionType, seqlen, keycard uint64) PokerSet2 {
	if actionType == Pass {
		return 0
	}
	var res PokerSet2
	cnt := PokerSet2(actionTypeBaseLenMap[actionType])
	var i uint64
	for i = 0; i < seqlen; i++ {
		res |= cnt << ((keycard + i) << 2)
	}
	return res
}

func (a Action) SetContext(c Action) Action {
	a[0] |= (c[0] << 32)
	return a
}

func (a Action) ContextActionType() uint64 {
	return (a[0] >> 32) & 255
}

func (a Action) ContextActionBrief() uint64 {
	return a[0] >> 32
}

func NewActionWithoutAff(actionType, seqlen, keycard uint64) Action {
	return Action{actionType | (seqlen << 8) | (keycard << 24),
		uint64(genPokerSet2WithActionStructionInfo(actionType, seqlen, keycard))}
}

func (a Action) SetAff(afflen int, ps PokerSet2) Action {
	a[0] |= uint64(afflen) << 16
	a[1] = uint64(ps.CombineWith(PokerSet2(a[1])))
	if a.ActionType() == Bomb && afflen > 0 {
		a[0] |= BW
	}
	return a
}

func NewActionWithPokerSet2(actionType, seqlen, afflen, keycard, ps2 uint64) Action {
	return Action{actionType | (seqlen << 8) | (afflen << 16) | (keycard << 24), ps2}
}

func (a Action) PokerCount() int {
	return actionTypeBaseLenMap[a[0]&255]*a.SeqLen() + a.AffLen()
	//return PokerSet2(a[1]).PokerCount()
}

func (a Action) String() string {
	if a == [2]uint64{} {
		return "无"
	}
	return actionTypeNameMap[a[0]&255] + "(" + PokerSet2(a[1]).String() + ")"
}

func (a Action) PokerSet2() PokerSet2 {
	return PokerSet2(a[1])
}

func (a Action) ActionType() uint64 {
	return a[0] & 255
}

func (a Action) SeqLen() int {
	return int((a[0] >> 8) & 255)
}

func (a Action) AffLen() int {
	return int((a[0] >> 16) & 255)
}

// 前24位
func (a Action) ActionStructure() uint64 {
	return a[0] & 0b111111111111111111111111
}

// 25到32位
func (a Action) KeyCard() uint64 {
	return a[0] >> 24 & 0xFF
}

// 前32位
func (a Action) WithoutContext() uint64 {
	return a[0] & 0b11111111111111111111111111111111
}

func (a Action) IsBomb() bool {
	return a[0]&255 == Bomb
}

func (a Action) PlayableInContext(c Action) bool {
	at := a.ActionType()
	ct := c.ActionType()
	if ct == Pass {
		return true
	} else {
		if at == Bomb {
			if ct == Bomb {
				if a.KeyCard() > c.KeyCard() {
					return true
				} else {
					return false
				}
			} else {
				return true
			}
		} else {
			if a.ActionStructure() == c.ActionStructure() {
				return a.WithoutContext() > c.WithoutContext()
			} else {
				return false
			}
		}
	}
}
