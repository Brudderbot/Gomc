package client

import (
	"encoding/json"

	"github.com/Brudderbot/Gomc/apis/buff"
	"github.com/Brudderbot/Gomc/impl/base"
	"github.com/Brudderbot/Gomc/impl/data/status"
)

// done

type PacketOResponse struct {
	Status status.Response
}

func (p *PacketOResponse) UUID() int32 {
	return 0x00
}

func (p *PacketOResponse) Push(writer buff.Buffer, conn base.Connection) {
	if text, err := json.Marshal(p.Status); err != nil {
		panic(err)
	} else {
		writer.PushTxt(string(text))
	}
}

type PacketOPong struct {
	Ping int64
}

func (p *PacketOPong) UUID() int32 {
	return 0x01
}

func (p *PacketOPong) Push(writer buff.Buffer, conn base.Connection) {
	writer.PushI64(p.Ping)
}
