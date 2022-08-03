package metaserver

import (
	"os"
	"time"
)

type Attr struct {
	Mtime         time.Time   // time of last modification
	Crtime        time.Time   // time of creation (OS X only)
	Mode          os.FileMode // file mode
	Uid           uint32      // owner uid
	Gid           uint32      // group gid
	Mime          string      // mime type
	TtlSec        int32       // ttl in seconds
	UserName      string
	GroupNames    []string
	SymlinkTarget string
	Md5           []byte
	FileSize      uint64
	Rdev          uint32
	Inode         uint64
}

func (attr Attr) IsDirectory() bool {
	return attr.Mode&os.ModeDir > 0
}

type Entry struct {
	// util.FullPath

	Attr
	Extended map[string][]byte

	// the following is for files
	// Chunks []*filer_pb.FileChunk `json:"chunks,omitempty"`

	// HardLinkId      HardLinkId
	// HardLinkCounter int32
	Content []byte
	// Remote  *filer_pb.RemoteEntry
	Quota int64
}

// func (entry *Entry) Size() uint64 {
// 	return maxUint64(maxUint64(TotalSize(entry.Chunks), entry.FileSize), uint64(len(entry.Content)))
// }

func (entry *Entry) Timestamp() time.Time {
	if entry.IsDirectory() {
		return entry.Crtime
	} else {
		return entry.Mtime
	}
}

func maxUint64(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}
