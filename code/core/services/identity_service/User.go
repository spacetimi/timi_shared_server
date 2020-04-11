package identity_service

import (
    "context"
    "errors"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service"
    "github.com/spacetimi/timi_shared_server/code/core/services/storage_service/storage_typedefs"
)

// Implements IBlob
type UserBlob struct {
    UserId int64
    CreatedTime int64
    LastLoginTime int64

    UserName string

    storage_typedefs.BlobDescriptor `bson:"ignore"`
}

func newUserBlob(userId int64) *UserBlob {
    user := UserBlob {
        UserId:userId,
    }
    user.BlobDescriptor = storage_typedefs.NewBlobDescriptor(storage_typedefs.STORAGE_SPACE_SHARED,
                                                             config.GetAppName() + "::user",
                                                             []string { "UserId" },
                                                             true)
    return &user
}

func loadUserBlobByUserId(userId int64, ctx context.Context) (*UserBlob, error) {
    user := newUserBlob(userId)

    err := storage_service.GetBlobByPrimaryKeys(user, ctx)
    if err != nil {
        return nil, errors.New("error getting user blob: " + err.Error())
    }

    return user, nil
}

