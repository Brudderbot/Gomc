package mode

import (
	"github.com/Brudderbot/Gomc/apis/util"
	"github.com/Brudderbot/Gomc/impl/base"
	"github.com/Brudderbot/Gomc/impl/prot/server"
)

/**
 * handshake
 */

func HandleState0(watcher util.Watcher) {

	watcher.SubAs(func(packet *server.PacketIHandshake, conn base.Connection) {
		conn.SetState(packet.State)
	})

}
