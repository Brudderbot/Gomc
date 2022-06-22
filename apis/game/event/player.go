package event

import "github.com/Brudderbot/Gomc/apis/ents"

type PlayerEvent struct {
	ents.Player
}

type PlayerJoinEvent struct {
	PlayerEvent
}

type PlayerQuitEvent struct {
	PlayerEvent
}
