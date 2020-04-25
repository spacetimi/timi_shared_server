package identity_service

import (
    "context"
    "errors"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service/storage_typedefs"
)

// Implements IBlob
type UserEmailToIdMappingBlob struct {
    UserEmailAddress string
    UserId int64

    storage_typedefs.BlobDescriptor `bson:"ignore"`
}

func newUserEmailToIdMapping(userEmailAddress string) *UserEmailToIdMappingBlob {
    ueidm := UserEmailToIdMappingBlob{
        UserEmailAddress:userEmailAddress,
    }
    ueidm.BlobDescriptor = storage_typedefs.NewBlobDescriptor(storage_typedefs.STORAGE_SPACE_SHARED,
                                                             config.GetAppName() + "::ueidm",
                                                             []string{ "UserEmailAddress" },
                                                             true)
    return &ueidm
}

func loadUserEmailToIdMappingByUserEmail(userEmailAddress string, ctx context.Context) (*UserEmailToIdMappingBlob, error) {
    ueidm := newUserEmailToIdMapping(userEmailAddress)

    err := storage_service.GetBlobByPrimaryKeys(ueidm, ctx)
    if err != nil {
        return nil, errors.New("error getting ueidm blob: " + err.Error())
    }

    return ueidm, nil
}

