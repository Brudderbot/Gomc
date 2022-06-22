package cmds

import (
	"github.com/Brudderbot/Gomc/apis/base"
	"github.com/Brudderbot/Gomc/apis/ents"
)

type Command interface {
	base.Named
	base.State

	Evaluate(sender ents.Sender, params []string)

	Complete(sender ents.Sender, params []string, output *[]string)
}
