package level

import "github.com/Brudderbot/Gomc/apis/base"

type Level interface {
	base.Named
	base.Unique

	Chunks() []Chunk

	GetChunk(x, z int) Chunk

	GetChunkIfLoaded(x, z int) Chunk

	GetBlock(x, y, z int) Block
}
