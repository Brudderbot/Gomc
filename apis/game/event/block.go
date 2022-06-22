package event

import (
	"github.com/Brudderbot/Gomc/apis/game/level"
)

type BlockEvent struct {
	level.Block
}

type BlockBreakEvent struct {
	BlockEvent
	PlayerEvent
	Cancellable
}
