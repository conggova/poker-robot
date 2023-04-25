package action

import (
	"math/bits"
	"strings"
)

func init() {
	bit := 32 << (^uint(0) >> 63)
	if bit != 64 {
		panic("Only suport 64 bit int.")
	}
}

// num to human readable string
/*
var numReadableMameMap = map[int]string{
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "10",
	11: "J",
	12: "Q",
	13: "K",
	15: "2",
	14: "A",
}
*/
// some poker cards , for the front 54 bits , each bit represents a card in 54 poker cards
// low bit to high bit , 3 to K , then A , then 2
type PokerSet uint64

// the card a bit represents
var bitCharMap = [54]byte{'3', '3', '3', '3', '4', '4', '4', '4', '5', '5', '5', '5', '6', '6', '6', '6',
	'7', '7', '7', '7', '8', '8', '8', '8', '9', '9', '9', '9', 'T', 'T', 'T', 'T',
	'J', 'J', 'J', 'J', 'Q', 'Q', 'Q', 'Q', 'K', 'K', 'K', 'K', 'A', 'A', 'A', 'A',
	'2', '2', '2', '2', 'z', 'Z'}

/*
var bitPowerMap = [54]int{0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5,
	6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 8, 9, 9, 9, 9, 10, 10, 10, 10, 11, 11, 11, 11, 12, 12, 12, 12, 13, 14}
*/
// some poker cards , dont distinguish 黑桃(spade)、红桃(heart)、梅花(club)、方块(dianmond)
// bit 0-3 represents poker 3 , 4-7 represents poker 4 , and so on
type PokerSet2 uint64

var powerCharMap = [15]byte{'3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K', 'A', '2', 'z', 'Z'}
var charPowerMap = map[byte]int{'3': 0, '4': 1, '5': 2, '6': 3, '7': 4, '8': 5, '9': 6, 'T': 7, 'J': 8, 'Q': 9, 'K': 10, 'A': 11, '2': 12, 'z': 13, 'Z': 14}

// information reduce
func (ps PokerSet) PokerSet2() PokerSet2 {
	var t PokerSet = 15
	var res PokerSet2
	var i uint
	for i < 15 {
		var a = bits.OnesCount(uint(t & ps))
		res |= PokerSet2(a << (i << 2)) //a left move i*4
		ps >>= 4                        //ps right move i*4
		i++
	}
	return res
}

func (ps PokerSet) PokerCount() int {
	return bits.OnesCount(uint(ps))
}

func (ps PokerSet2) PokerCount() int {
	var sum PokerSet2
	var t PokerSet2 = 15 //1111
	for i := 0; i < 15; i++ {
		sum += t & ps
		ps >>= 4
	}
	return int(sum)
}

func (ps PokerSet) CombineWith(ps1 PokerSet) PokerSet {
	return ps | ps1
}

func (ps PokerSet2) CombineWith(ps1 PokerSet2) PokerSet2 {
	return ps + ps1
}

func (ps PokerSet) Subtract(ps1 PokerSet) PokerSet {
	return ps & (^(ps & ps1))
}

// without check
func (ps PokerSet2) Subtract(ps1 PokerSet2) PokerSet2 {
	/* for debug
	if !ps.Covers(ps1) {
		panic(fmt.Sprintf("subtract ps %s , ps1 %s", ps, ps1))
	}
	*/
	return ps - ps1
}

func (ps PokerSet) String() string {
	charArray := make([]byte, 0, 20)
	var t PokerSet = 1
	for i := 0; i < 54; i++ {
		if t&ps != 0 {
			charArray = append(charArray, bitCharMap[i])
		}
		t <<= 1
	}
	return string(charArray)
}

func (ps PokerSet2) String() string {
	charArray := make([]byte, 0, 20)
	var t PokerSet2 = 15 //1111
	for i := 0; i < 15; i++ {
		cnt := t & ps
		ps >>= 4
		c := powerCharMap[i]
		for j := 0; j < int(cnt); j++ {
			charArray = append(charArray, c)
		}
	}
	return string(charArray)
}

func (ps PokerSet) MaxPowerPoker() int {
	var t PokerSet = 0b111100000000000000000000000000000000000000000000000000000000
	for i := 14; i >= 0; i-- {
		if t&ps > 0 {
			return i
		}
		t >>= 4
	}
	return -1
}

func (ps PokerSet2) MaxPowerPoker() int {
	var t PokerSet2 = 0b111100000000000000000000000000000000000000000000000000000000
	for i := 14; i >= 0; i-- {
		if t&ps > 0 {
			return i
		}
		t >>= 4
	}
	return -1
}

func (ps PokerSet2) HaveNonSingle() bool {
	var t PokerSet2 = 0b1111
	var seqlen int
	for i := 0; i < 12; i++ {
		cnt := t & ps
		ps >>= 4
		if cnt > 0 {
			seqlen += 1
			if seqlen >= 5 {
				return true
			}
			if cnt > 1 {
				return true
			}
		} else {
			seqlen = 0
		}
	}
	return false
}

func ParsePokerSet2(ipt string) (bool, PokerSet2) {
	ipt = strings.ToUpper(ipt)
	byteSlice := []byte(ipt)
	powerCnt := map[int]uint64{}
	for _, char := range byteSlice {
		if v, ok := charPowerMap[char]; ok {
			powerCnt[v] += 1
		} else {
			return false, 0
		}
	}
	var res PokerSet2
	for i := 0; i < 15; i++ {
		res |= PokerSet2(powerCnt[i] << (i << 2))
	}
	return true, res
}

func (ps PokerSet2) Covers(ps1 PokerSet2) bool {
	if ps >= ps1 && (ps-ps1)&0x888888888888888 == 0 { //不发生借位
		return true
	}
	return false
}

func (ps PokerSet2) parse() (pkCnt, sCnt, cCnt, tCnt, bCnt int, sList, cList, tList [][2]uint64, bList []uint64) {
	var t PokerSet2 = 0b1111
	var start1, start2, start3, i uint64 = 1000, 1000, 1000, 0
	for i = 0; i < 15; i++ {
		cnt := int(t & ps)
		pkCnt += cnt
		ps >>= 4
		if cnt == 0 || i > 11 { //11 is A
			if start1 != 1000 {
				t := i - start1
				sCnt += int(t)
				sList = append(sList, [2]uint64{start1, t})
				start1 = 1000
			}
			if start2 != 1000 {
				t := i - start2
				cCnt += int(t)
				cList = append(cList, [2]uint64{start2, t})
				start2 = 1000
			}
			if start3 != 1000 {
				t := i - start3
				tCnt += int(t)
				tList = append(tList, [2]uint64{start3, t})
				start3 = 1000
			}
		}
		if cnt >= 1 {
			if start1 == 1000 {
				start1 = i
			}
			if cnt == 1 {
				if start2 != 1000 {
					t := i - start2
					cCnt += int(t)
					cList = append(cList, [2]uint64{start2, t})
					start2 = 1000
				}
				if start3 != 1000 {
					t := i - start3
					tCnt += int(t)
					tList = append(tList, [2]uint64{start3, t})
					start3 = 1000
				}
			} else if cnt >= 2 {
				if start2 == 1000 {
					start2 = i
				}
				if cnt == 2 {
					if start3 != 1000 {
						t := i - start3
						tCnt += int(t)
						tList = append(tList, [2]uint64{start3, t})
						start3 = 1000
					}
				} else if cnt >= 3 {
					if start3 == 1000 {
						start3 = i
					}
					if cnt == 4 {
						bList = append(bList, i)
						bCnt += 1
					}
				}
			}
		}
	}
	return
}
