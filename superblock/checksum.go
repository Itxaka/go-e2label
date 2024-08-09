package superblock

import (
	"hash/crc32"
)

func NewChecksum() *Checksum {
	return &Checksum{
		val:   0,
		table: crc32.MakeTable(crc32.Castagnoli),
	}
}

type Checksum struct {
	val   uint32
	table *crc32.Table
}

func (cs *Checksum) Write(b []byte) (n int, err error) {
	cs.val = crc32.Update(cs.val, cs.table, b)
	return len(b), nil
}

func (cs *Checksum) Get() uint32 {
	return ^cs.val
}
