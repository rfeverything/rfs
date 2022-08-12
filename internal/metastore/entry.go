package metastore

import (
	"github.com/google/uuid"
)

type MetaEntry struct {
	File  string `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	Size  uint32 `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Owner string `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	Group string `protobuf:"bytes,4,opt,name=group,proto3" json:"group,omitempty"`
	//Optional fields for file.
	BlockReplication uint32        `protobuf:"varint,5,opt,name=block_replication,json=blockReplication,proto3" json:"block_replication,omitempty"`
	BlockSize        uint32        `protobuf:"varint,6,opt,name=block_size,json=blockSize,proto3" json:"block_size,omitempty"`
	Locations        *LocatedBlock `protobuf:"bytes,7,opt,name=locations,proto3" json:"locations,omitempty"`
	FileId           uint64        `protobuf:"varint,8,opt,name=fileId,proto3" json:"fileId,omitempty"`
	Flags            uint32        `protobuf:"varint,9,opt,name=flags,proto3" json:"flags,omitempty"`
	Namespace        string        `protobuf:"bytes,10,opt,name=namespace,proto3" json:"namespace,omitempty"`

	Fid     uuid.UUID
	Volumes []string
}

type LocatedBlock struct {
	Block     *Block              `protobuf:"bytes,1,opt,name=block,proto3" json:"block,omitempty"`
	Offset    uint64              `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Locations []*VolumeServerInfo `protobuf:"bytes,3,rep,name=locations,proto3" json:"locations,omitempty"`
}

type Block struct {
	BlockId uint64 `protobuf:"varint,1,opt,name=blockId,proto3" json:"blockId,omitempty"`
	Size    uint64 `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
}

type VolumeServerInfo struct {
	Host     string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port     uint32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	VolumeId uint32 `protobuf:"varint,3,opt,name=volumeId,proto3" json:"volumeId,omitempty"`
}
