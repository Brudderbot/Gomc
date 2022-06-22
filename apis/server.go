package apis

import (
	"sync"

	"github.com/Brudderbot/Gomc/apis/cmds"
	"github.com/Brudderbot/Gomc/apis/ents"
	"github.com/Brudderbot/Gomc/apis/logs"
	"github.com/Brudderbot/Gomc/apis/task"
	"github.com/Brudderbot/Gomc/apis/util"
	"github.com/Brudderbot/Gomc/apis/uuid"

	apis_base "github.com/Brudderbot/Gomc/apis/base"
	impl_base "github.com/Brudderbot/Gomc/impl/base"
)

type Server interface {
	apis_base.State

	Logging() *logs.Logging

	Command() *cmds.CommandManager

	Tasking() *task.Tasking

	Watcher() util.Watcher

	Players() []ents.Player

	ForEachPlayer(cb func(ents.Player))

	ConnByUUID(uuid uuid.UUID) impl_base.Connection

	PlayerByUUID(uuid uuid.UUID) ents.Player

	PlayerByConn(conn impl_base.Connection) ents.Player

	ServerVersion() string

	Broadcast(message string)
}

var instance *Server
var syncOnce sync.Once

func MinecraftServer() Server {
	if instance == nil {
		panic("server is unavailable")
	}

	return *instance
}

func SetMinecraftServer(server Server) {
	syncOnce.Do(func() {
		instance = &server
	})
}
