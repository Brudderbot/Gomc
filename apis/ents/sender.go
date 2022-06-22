package ents

import (
	"github.com/Brudderbot/Gomc/apis/base"
)

type Sender interface {
	base.Named

	SendMessage(message ...interface{})
}
