package mode

import (
	"fmt"
	"time"

	"github.com/Brudderbot/Gomc/apis"
	"github.com/Brudderbot/Gomc/apis/data"
	"github.com/Brudderbot/Gomc/apis/data/chat"
	"github.com/Brudderbot/Gomc/apis/data/msgs"
	"github.com/Brudderbot/Gomc/apis/ents"
	"github.com/Brudderbot/Gomc/apis/game"
	"github.com/Brudderbot/Gomc/apis/logs"
	"github.com/Brudderbot/Gomc/apis/task"
	"github.com/Brudderbot/Gomc/apis/util"
	"github.com/Brudderbot/Gomc/impl/base"
	"github.com/Brudderbot/Gomc/impl/data/client"
	"github.com/Brudderbot/Gomc/impl/data/plugin"
	"github.com/Brudderbot/Gomc/impl/data/values"
	impl_level "github.com/Brudderbot/Gomc/impl/game/level"

	impl_event "github.com/Brudderbot/Gomc/impl/game/event"

	client_packet "github.com/Brudderbot/Gomc/impl/prot/client"
	server_packet "github.com/Brudderbot/Gomc/impl/prot/server"
)

func HandleState3(watcher util.Watcher, logger *logs.Logging, tasking *task.Tasking, join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection) {
	tasking.EveryTime(10*time.Second, func(*task.Task) {
		api := apis.MinecraftServer()

		api.ForEachPlayer(func(player ents.Player) {
			// also probably add one that returns both the player and their connection
			conn := api.ConnByUUID(player.UUID())

			// keep player connection alive via keep alive
			conn.SendPacket(&client_packet.PacketOKeepAlive{KeepAliveID: time.Now().UnixNano() / 1e6})
		})
	})

	watcher.SubAs(func(packet *server_packet.PacketIKeepAlive, conn base.Connection) {
		logger.DataF("player %s is being kept alive", conn.Address())
	})

	watcher.SubAs(func(packet *server_packet.PacketIPluginMessage, conn base.Connection) {
		api := apis.MinecraftServer()

		player := api.PlayerByConn(conn)
		if player == nil {
			return // log no player found?
		}

		api.Watcher().PubAs(impl_event.PlayerPluginMessagePullEvent{
			Conn: base.PlayerAndConnection{
				Connection: conn,
				Player:     player,
			},
			Channel: packet.Message.Chan(),
			Message: packet.Message,
		})
	})

	watcher.SubAs(func(packet *server_packet.PacketIChatMessage, conn base.Connection) {
		api := apis.MinecraftServer()

		who := api.PlayerByConn(conn)
		out := msgs.
			New(who.Name()).SetColor(chat.White).
			Add(":").SetColor(chat.Gray).
			Add(" ").
			Add(chat.Translate(packet.Message)).SetColor(chat.White).
			AsText() // why not just use translate?

		api.Broadcast(out)
	})

	watcher.SubAs(func(packet *server_packet.PacketIPlayerAnimation, conn base.Connection) {
		api := apis.MinecraftServer()

		who := api.PlayerByConn(conn)
		outPacket := &client_packet.PacketOEntityAnimation{
			Entity:    who,
			Animation: packet.Hand,
		}

		api.ForEachPlayer(func(player ents.Player) {
			if player == who {
				return
			}

			// also probably add one that returns both the player and their connection
			conn := api.ConnByUUID(player.UUID())

			// keep player connection alive via keep alive
			conn.SendPacket(outPacket)
		})
	})

	go func() {
		api := apis.MinecraftServer()

		for conn := range join {
			thisPlayerSpawnPacket := &client_packet.PacketOSpawnPlayer{
				EntityID:   conn.Player.EntityUUID(),
				PlayerUUID: conn.UUID(),
				Location: data.PositionF{
					X: 0,
					Y: 10,
					Z: 0,
				},
			}
			api.ForEachPlayer(func(player ents.Player) {
				if player == conn.Player {
					return
				}

				// also probably add one that returns both the player and their connection
				conn := api.ConnByUUID(player.UUID())

				// keep player connection alive via keep alive
				conn.SendPacket(thisPlayerSpawnPacket)
			})

			apis.MinecraftServer().Watcher().PubAs(impl_event.PlayerConnJoinEvent{Conn: conn})

			conn.SendPacket(&client_packet.PacketOJoinGame{
				EntityID:      int32(conn.EntityUUID()),
				Hardcore:      false,
				GameMode:      game.CREATIVE,
				Dimension:     game.OVERWORLD,
				HashedSeed:    values.DefaultWorldHashedSeed,
				MaxPlayers:    10,
				LevelType:     game.DEFAULT,
				ViewDistance:  12,
				ReduceDebug:   false,
				RespawnScreen: false,
			})

			conn.SendPacket(&client_packet.PacketOPluginMessage{
				Message: &plugin.Brand{
					Name: chat.Translate(fmt.Sprintf("&b%s&r &a%s&r", "GoLangMc", apis.MinecraftServer().ServerVersion())),
				},
			})

			conn.SendPacket(&client_packet.PacketOServerDifficulty{
				Difficulty: game.PEACEFUL,
				Locked:     true,
			})

			conn.SendPacket(&client_packet.PacketOPlayerAbilities{
				Abilities: client.PlayerAbilities{
					Invulnerable: true,
					Flying:       true,
					AllowFlight:  true,
					InstantBuild: false,
				},
				FlyingSpeed: 0.05, // default value
				FieldOfView: 0.1,  // default value
			})

			conn.SendPacket(&client_packet.PacketOHeldItemChange{
				Slot: client.SLOT_0,
			})

			conn.SendPacket(&client_packet.PacketODeclareRecipes{})

			conn.SendPacket(&client_packet.PacketOPlayerLocation{
				ID: 0,
				Location: data.Location{
					PositionF: data.PositionF{
						X: 0,
						Y: 10,
						Z: 0,
					},
					RotationF: data.RotationF{
						AxisX: 0,
						AxisY: 0,
					},
				},
				Relative: client.Relativity{},
			})

			conn.SendPacket(&client_packet.PacketOPlayerInfo{
				Action: client.AddPlayer,
				Values: []client.PlayerInfo{
					&client.PlayerInfoAddPlayer{Player: conn.Player},
				},
			})

			conn.SendPacket(&client_packet.PacketOEntityMetadata{Entity: conn.Player})

			level := impl_level.NewLevel("test")
			impl_level.GenSuperFlat(level, 6)

			for _, chunk := range level.Chunks() {
				conn.SendPacket(&client_packet.PacketOChunkData{Chunk: chunk})
			}

			logger.DataF("chunks sent to player: %s", conn.Player.Name())

			conn.SendPacket(&client_packet.PacketOPlayerLocation{
				ID: 1,
				Location: data.Location{
					PositionF: data.PositionF{
						X: 0,
						Y: 10,
						Z: 0,
					},
					RotationF: data.RotationF{
						AxisX: 0,
						AxisY: 0,
					},
				},
				Relative: client.Relativity{},
			})
		}
	}()

	go func() {
		for conn := range quit {
			apis.MinecraftServer().Watcher().PubAs(impl_event.PlayerConnQuitEvent{Conn: conn})
		}
	}()
}
