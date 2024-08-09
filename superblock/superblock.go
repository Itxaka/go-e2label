package superblock

import (
	"github.com/lunixbochs/struc"
	"io"
	"os"
	"syscall"
)

const BlockStart = 1024

type Superblock struct {
	InodeCount           uint32     `struc:"uint32,little"`
	BlockcountLo         uint32     `struc:"uint32,little"`
	RBlockcountLo        uint32     `struc:"uint32,little"`
	FreeBlockcountLo     uint32     `struc:"uint32,little"`
	FreeInodecount       uint32     `struc:"uint32,little"`
	FirstDataBlock       uint32     `struc:"uint32,little"`
	LogBlockSize         uint32     `struc:"uint32,little"`
	LogClusterSize       uint32     `struc:"uint32,little"`
	BlockperGroup        uint32     `struc:"uint32,little"`
	ClusterperGroup      uint32     `struc:"uint32,little"`
	InodeperGroup        uint32     `struc:"uint32,little"`
	Mtime                uint32     `struc:"uint32,little"`
	Wtime                uint32     `struc:"uint32,little"`
	MntCount             uint16     `struc:"uint16,little"`
	MaxMntCount          uint16     `struc:"uint16,little"`
	Magic                uint16     `struc:"uint16,little"`
	State                uint16     `struc:"uint16,little"`
	Errors               uint16     `struc:"uint16,little"`
	MinorRevLevel        uint16     `struc:"uint16,little"`
	Lastcheck            uint32     `struc:"uint32,little"`
	Checkinterval        uint32     `struc:"uint32,little"`
	CreatorOs            uint32     `struc:"uint32,little"`
	RevLevel             uint32     `struc:"uint32,little"`
	DefResuid            uint16     `struc:"uint16,little"`
	DefResgid            uint16     `struc:"uint16,little"`
	FirstIno             uint32     `struc:"uint32,little"`
	InodeSize            uint16     `struc:"uint16,little"`
	BlockGroupNr         uint16     `struc:"uint16,little"`
	FeatureCompat        uint32     `struc:"uint32,little"`
	FeatureIncompat      uint32     `struc:"uint32,little"`
	FeatureRoCompat      uint32     `struc:"uint32,little"`
	Uuid                 [16]byte   `struc:"[16]byte"`
	VolumeName           [16]byte   `struc:"[16]byte"`
	LastMounted          [64]byte   `struc:"[64]byte"`
	AlgorithmUsageBitmap uint32     `struc:"uint32,little"`
	PreallocBlocks       byte       `struc:"byte"`
	PreallocDirBlocks    byte       `struc:"byte"`
	ReservedGdtBlocks    uint16     `struc:"uint16,little"`
	JournalUuid          [16]byte   `struc:"[16]byte"`
	JournalInum          uint32     `struc:"uint32,little"`
	JournalDev           uint32     `struc:"uint32,little"`
	LastOrphan           uint32     `struc:"uint32,little"`
	HashSeed             [4]uint32  `struc:"[4]uint32,little"`
	DefHashVersion       byte       `struc:"byte"`
	JnlBackupType        byte       `struc:"byte"`
	DescSize             uint16     `struc:"uint16,little"`
	DefaultMountOpts     uint32     `struc:"uint32,little"`
	FirstMetaBg          uint32     `struc:"uint32,little"`
	MkfTime              uint32     `struc:"uint32,little"`
	JnlBlocks            [17]uint32 `struc:"[17]uint32,little"`
	BlockcountHi         uint32     `struc:"uint32,little"`
	RBlockcountHi        uint32     `struc:"uint32,little"`
	FreeBlockcountHi     uint32     `struc:"uint32,little"`
	MinExtraIsize        uint16     `struc:"uint16,little"`
	WantExtraIsize       uint16     `struc:"uint16,little"`
	Flags                uint32     `struc:"uint32,little"`
	RaidStride           uint16     `struc:"uint16,little"`
	MmpUpdateInterval    uint16     `struc:"uint16,little"`
	MmpBlock             uint64     `struc:"uint64,little"`
	RaidStripeWidth      uint32     `struc:"uint32,little"`
	LogGroupperFlex      byte       `struc:"byte"`
	ChecksumType         byte       `struc:"byte"`
	EncryptionLevel      byte       `struc:"byte"`
	ReservedPad          byte       `struc:"byte"`
	KbyteWritten         uint64     `struc:"uint64,little"`
	SnapshotInum         uint32     `struc:"uint32,little"`
	SnapshotId           uint32     `struc:"uint32,little"`
	SnapshotRBlockcount  uint64     `struc:"uint64,little"`
	SnapshotList         uint32     `struc:"uint32,little"`
	ErrorCount           uint32     `struc:"uint32,little"`
	FirstErrorTime       uint32     `struc:"uint32,little"`
	FirstErrorIno        uint32     `struc:"uint32,little"`
	FirstErrorBlock      uint64     `struc:"uint64,little"`
	FirstErrorFunc       [32]byte   `struc:"[32]pad"`
	FirstErrorLine       uint32     `struc:"uint32,little"`
	LastErrorTime        uint32     `struc:"uint32,little"`
	LastErrorIno         uint32     `struc:"uint32,little"`
	LastErrorLine        uint32     `struc:"uint32,little"`
	LastErrorBlock       uint64     `struc:"uint64,little"`
	LastErrorFunc        [32]byte   `struc:"[32]pad"`
	MountOpts            [64]byte   `struc:"[64]pad"`
	UsrQuotaInum         uint32     `struc:"uint32,little"`
	GrpQuotaInum         uint32     `struc:"uint32,little"`
	OverheadClusters     uint32     `struc:"uint32,little"`
	BackupBgs            [2]uint32  `struc:"[2]uint32,little"`
	EncryptAlgos         [4]byte    `struc:"[4]pad"`
	EncryptPwSalt        [16]byte   `struc:"[16]pad"`
	LpfIno               uint32     `struc:"uint32,little"`
	PrjQuotaInum         uint32     `struc:"uint32,little"`
	ChecksumSeed         uint32     `struc:"uint32,little"`
	Reserved             [98]uint32 `struc:"[98]uint32,little"`
	Checksum             uint32     `struc:"uint32,little"`
}

func (sb *Superblock) calculateChecksum() uint32 {
	cs := NewChecksum()
	size, _ := struc.Sizeof(sb)
	struc.Pack(LimitWriter(cs, int64(size)-4), sb)
	sb.Checksum = cs.Get()
	return cs.Get()
}

func (sb *Superblock) CalculateNewChecksum() uint32 {
	return sb.calculateChecksum()
}

func (sb *Superblock) CalculateNewChecksumAndWriteIt(devicePath string) error {
	f, err := os.OpenFile(devicePath, syscall.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	sb.Checksum = sb.calculateChecksum()

	_, err = f.Seek(BlockStart, 0)
	if err != nil {
		return err
	}
	return struc.Pack(f, sb)
}

// GetSuperBlock accepts a path to a device (file or device) and tries to get the SuperBlock from it
func GetSuperBlock(devicePath string) (*Superblock, error) {
	ret := &Superblock{}
	f, err := os.OpenFile(devicePath, syscall.O_RDWR, 0755)
	if err != nil {
		return ret, err
	}
	defer f.Close()

	_, err = f.Seek(BlockStart, 0)
	if err != nil {
		return ret, err
	}

	err = struc.Unpack(f, ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// A limitedWriter writes to W but limits the amount of
// data written to just N bytes. Each call to Write
// updates N to reflect the new amount remaining.
// Write returns EOF when N <= 0 or when the underlying W returns EOF.
type limitedWriter struct {
	W io.Writer
	N int64
}

func LimitWriter(w io.Writer, n int64) io.Writer { return &limitedWriter{w, n} }

func (lw *limitedWriter) Write(p []byte) (n int, err error) {
	if lw.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > lw.N {
		p = p[0:lw.N]
	}
	n, err = lw.W.Write(p)
	lw.N -= int64(n)
	return
}
