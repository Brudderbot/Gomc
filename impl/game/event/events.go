package event

import (
	"github.com/Brudderbot/Gomc/impl/base"
	"github.com/Brudderbot/Gomc/impl/data/plugin"
)

type PlayerConnJoinEvent struct {
	Conn base.PlayerAndConnection
}

type PlayerConnQuitEvent struct {
	Conn base.PlayerAndConnection
}

type PlayerPluginMessagePullEvent struct {
	Conn base.PlayerAndConnection

	Channel string
	Message plugin.Message
}
