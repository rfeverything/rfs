package metaserver

type ChunkView struct {
	FileId      string
	Offset      int64
	Size        uint64
	LogicOffset int64 // actual offset in the file, for the data specified via [offset, offset+size) in current chunk
	ChunkSize   uint64
	IsGzipped   bool
}
