package base

import "github.com/Brudderbot/Gomc/apis/uuid"

type Unique interface {
	UUID() uuid.UUID
}
