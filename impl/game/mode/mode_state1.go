package mode

import (
	"github.com/Brudderbot/Gomc/apis/util"
	"github.com/Brudderbot/Gomc/impl/base"
	"github.com/Brudderbot/Gomc/impl/data/status"
	"github.com/Brudderbot/Gomc/impl/prot/client"
	"github.com/Brudderbot/Gomc/impl/prot/server"
)

/**
 * status
 */

func HandleState1(watcher util.Watcher) {

	watcher.SubAs(func(packet *server.PacketIRequest, conn base.Connection) {
		response := client.PacketOResponse{Status: status.DefaultResponse()}
		conn.SendPacket(&response)
	})

	watcher.SubAs(func(packet *server.PacketIPing, conn base.Connection) {
		response := client.PacketOPong{Ping: packet.Ping}
		conn.SendPacket(&response)
	})

}
