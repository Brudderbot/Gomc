package ents

import "github.com/Brudderbot/Gomc/apis/base"

type Entity interface {
	Sender
	base.Unique

	EntityUUID() int64
}
