package storage_typedefs

type StorageSpace int

const (
	STORAGE_SPACE_SHARED StorageSpace = iota
	STORAGE_SPACE_APP
)

type IBlob interface {
	GetStorageSpace() StorageSpace
	GetBlobName() string
	GetPrimaryKeys() []string
	GetVersion() int
	IsRedisAllowed() bool
}

/***** Concrete Types ***********************************************************/

type BlobDescriptor struct { // Implements IBlob
	space          StorageSpace
	blobName       string
	primaryKeys    []string
	version        int
	isRedisAllowed bool
}

func NewBlobDescriptor(space StorageSpace, blobName string, primaryKeys []string, version int, isRedisAllowed bool) BlobDescriptor {
	bd := BlobDescriptor{
		space:          space,
		blobName:       blobName,
		primaryKeys:    primaryKeys,
		isRedisAllowed: isRedisAllowed,
	}
	return bd
}

func (bd *BlobDescriptor) GetStorageSpace() StorageSpace {
	return bd.space
}

func (bd *BlobDescriptor) GetBlobName() string {
	return bd.blobName
}

func (bd *BlobDescriptor) GetPrimaryKeys() []string {
	return bd.primaryKeys
}

func (bd *BlobDescriptor) GetVersion() int {
	return bd.version
}

func (bd *BlobDescriptor) IsRedisAllowed() bool {
	return bd.isRedisAllowed
}
