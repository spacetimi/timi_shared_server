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
}

/***** Concrete Types ***********************************************************/

type BlobDescriptor struct {        // Implements IBlob
    space          StorageSpace
    blobName       string
    primaryKeys    []string
}

func NewBlobDescriptor(space StorageSpace, blobName string, primaryKeys []string) BlobDescriptor {
    bd := BlobDescriptor{
        space:space,
        blobName:blobName,
        primaryKeys:primaryKeys,
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
