package strategy

import (
	"github.com/conggova/poker-robot/pkg/action"
)

type TCheatFlag int
type TCheatMethod int

const (
	CommInSecret  TCheatMethod = iota //私下通信，互知底牌
	ShareInterest                     //利益共通
)

const (
	NoCheat TCheatFlag = iota
	WithPreer
	WithNexter
)

type PlayContext struct {
	PreerAction          action.Action
	NexterAction         action.Action
	PreerPkCnt           int
	NexterPkCnt          int
	GameOpenCard         bool
	OthersRemainPokerSet action.PokerSet2 //必有
	PreerRemainPokerSet  action.PokerSet2 //GameOpenCard为true才有有效值
	NexterRemainPokerSet action.PokerSet2 //GameOpenCard为true才有有效值
	PlayLog              []action.Action
	PreerPlayLog         []action.Action
	NexterPlayLog        []action.Action
}

type StrategyContext struct {
	RemainPokerSet action.PokerSet2 //
	CheatFlag      TCheatFlag
	CheatMethod    TCheatMethod
	OpenCard4Me    bool //为True时PreerRemainPokerSet和NexterRemainPokerSet必有有效值
	PlayContext
}

type GameBuildInfo struct {
	Strategys   [3]Strategy
	OpenCard    bool
	Partners    [2]int
	CheatMethod TCheatMethod
}

type GameRestoreInfo struct {
	LastPlayer      int                 //最后是谁出的牌
	RemainPokerSets [3]action.PokerSet2 //各家余牌
	PlayLogs        [3][]action.Action  //各家出牌记录
}

// 为不让Strategy不直接依赖Game（会出现循环依赖），让Game实现Simulator接口
type Simulator interface {
	BuildGame(GameBuildInfo) Simulator //通过定义在接口里的构造函数，其实实现了循环依赖
	RestoreGame(GameRestoreInfo)
	RunRestoredGame()
	GetProfits() [3]int
}

type Strategy interface {
	MakeDecision(c StrategyContext) action.Action
}
