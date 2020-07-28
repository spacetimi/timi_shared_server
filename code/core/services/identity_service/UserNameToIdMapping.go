package identity_service

import (
	"context"
	"errors"

	"github.com/spacetimi/timi_shared_server/v2/code/config"
	"github.com/spacetimi/timi_shared_server/v2/code/core/services/storage_service"
	"github.com/spacetimi/timi_shared_server/v2/code/core/services/storage_service/storage_typedefs"
)

const kUNIDMVersion = 1

// Implements IBlob
type UserNameToIdMappingBlob struct {
	UserName     string
	PasswordHash string
	UserId       int64

	storage_typedefs.BlobDescriptor `bson:"ignore"`
}

func newUserNameToIdMapping(userName string) *UserNameToIdMappingBlob {
	uidm := UserNameToIdMappingBlob{
		UserName: userName,
	}
	uidm.BlobDescriptor = storage_typedefs.NewBlobDescriptor(storage_typedefs.STORAGE_SPACE_SHARED,
		config.GetAppName()+"::uidm",
		[]string{"UserName"},
		kUNIDMVersion,
		true)
	return &uidm
}

func loadUserNameToIdMappingByUserName(userName string, ctx context.Context) (*UserNameToIdMappingBlob, error) {
	uidm := newUserNameToIdMapping(userName)

	err := storage_service.GetBlobByPrimaryKeys(uidm, ctx)
	if err != nil {
		return nil, errors.New("error getting uidm blob: " + err.Error())
	}

	return uidm, nil
}
