package metaserver

import (
	"errors"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	rfspb "github.com/rfeverything/rfs/internal/proto/rfs"
)

type Entry struct {
	Path string

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

func (a *Attr) IsDir() bool {
	return a.Mode.IsDir()
}

func (e *Entry) ToExistingProtoEntry(message *rfspb.Entry) error {
	if e == nil {
		return errors.New("entry is nil")
	}
	message.Name = e.Path
	message.Attributes = &rfspb.FuseAttributes{
		Mtime:    e.Mtime.Unix(),
		Crtime:   e.Crtime.Unix(),
		FileMode: uint32(e.Mode),
		Uid:      e.Uid,
		Gid:      e.Gid,
		Mime:     e.Mime,
		Md5:      e.Md5,
		FileSize: e.FileSize,
		Rdev:     e.Rdev,
		Inode:    e.Inode,
	}
	message.Extended = e.Extended
	message.Chunks = e.Chunks
	return nil
}

func (e *Entry) EncodeAttributesAndChunks() ([]byte, error) {
	msg := &rfspb.Entry{}
	e.ToExistingProtoEntry(msg)
	return proto.Marshal(msg)
}

func FromProtoEntry(msg *rfspb.Entry, e *Entry) {
	if msg == nil {
		return
	}
	e.Path = msg.Name
	e.Mtime = time.Unix(msg.Attributes.Mtime, 0)
	e.Crtime = time.Unix(msg.Attributes.Crtime, 0)
	e.Mode = os.FileMode(msg.Attributes.FileMode)
	e.Uid = msg.Attributes.Uid
	e.Gid = msg.Attributes.Gid
	e.Mime = msg.Attributes.Mime
	e.Md5 = msg.Attributes.Md5
	e.FileSize = msg.Attributes.FileSize
	e.Rdev = msg.Attributes.Rdev
	e.Inode = msg.Attributes.Inode
	e.Extended = msg.Extended
	e.Chunks = msg.Chunks
	return
}

func (e *Entry) DecodeAttributesAndChunks(data []byte) error {
	msg := &rfspb.Entry{}
	if err := proto.Unmarshal(data, msg); err != nil {
		return err
	}
	FromProtoEntry(msg, e)
	return nil
}
