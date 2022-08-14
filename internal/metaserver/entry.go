package metaserver

import (
	"os"
	"time"

	rfspb "github.com/rfeverything/rfs/internal/proto/rfs"
)

type Entry struct {
	Dir string

	Attr
	Extended map[string][]byte

	// the following is for files
	Chunks []*rfspb.FileChunk `json:"chunks,omitempty"`

	// HardLinkId      HardLinkId
	// HardLinkCounter int32
	Content []byte

	Quota int64
}

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
